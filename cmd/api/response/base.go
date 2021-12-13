package response

import "github.com/gin-gonic/gin"

const (
	Success          = "success"
	ParamInvalid     = "param_invalid"
	TokenInvalid     = "token_invalid"
	PermissionDenied = "permission_denied"
	UserNotFound     = "user_not_found"
	TaskNotFound     = "task_not_found"
)

type ResWithRemark struct {
	Response
	Remark `json:"remark"`
}

type Remark []interface{}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func MkResponse(ctx *gin.Context, code int, msg string, data interface{}, remark ...interface{}) {
	ctx.JSON(code, gin.H{
		"code":   code,
		"msg":    msg,
		"data":   data,
		"remark": remark,
	})
}
