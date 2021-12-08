package operation

import (
	"bytes"
	"encoding/json"
	"net/url"
	"time"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

const logReq = "logReq"

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

type ResWithRemark struct {
	Res
	Remark `json:"remark"`
}

type Res struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Remark []interface{}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	bWithoutRemark := removeRemarkFromResponse(b)
	return w.ResponseWriter.Write(bWithoutRemark)
}

func removeRemarkFromResponse(b []byte) []byte {
	res := Res{}
	_ = jsoniter.Unmarshal(b, &res)
	resB, _ := jsoniter.Marshal(&res)
	return resB
}

func Log() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		ctx.Next()

		response := blw.body.String()

		user := helper.GetUserClaims(ctx)
		if user == nil {
			return
		}

		now := time.Now()

		clients.WriteDBCli.Create(&model.OperationLog{
			Base: model.Base{
				CreateAt: &now,
				UpdateAt: &now,
			},
			Handler:  getHandlerFunc(ctx),
			Params:   getParams(ctx),
			Info:     getReq(ctx),
			Operator: user.UserId,
			Response: response,
		})
	}
}

func getHandlerFunc(ctx *gin.Context) string {
	handlers := ctx.HandlerNames()
	handler := handlers[len(handlers)-1]
	return utils.GetStringSuffix(handler)
}

func getParams(ctx *gin.Context) string {
	params, _ := url.ParseQuery(ctx.Request.URL.RawQuery)
	p, _ := jsoniter.Marshal(params)
	return string(p)
}

func LogReq(ctx *gin.Context, v interface{}) {
	ctx.Set(logReq, v)
}

func getReq(ctx *gin.Context) string {
	v, _ := ctx.Get(logReq)
	info, _ := json.Marshal(v)
	return string(info)
}