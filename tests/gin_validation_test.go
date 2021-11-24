package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/galaxy-future/BridgX/cmd/api/middleware"
	"github.com/galaxy-future/BridgX/cmd/api/middleware/validation"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type validateCase struct {
	Name string `form:"name" binding:"required,charTypeGT3" json:"name"`
}

func testServer(t *testing.T) *gin.Engine {
	middleware.Init()
	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := validation.RegisterValidators(v)
		if err != nil {
			t.Errorf("register tools failed.err:[%s]", err.Error())
			return r
		}
	}
	return r
}

func TestValidator_GET(t *testing.T) {
	r := testServer(t)
	r.GET("/ping", func(ctx *gin.Context) {
		c := validateCase{}
		err := ctx.ShouldBindWith(&c, binding.Query)
		if err == nil {
			t.Errorf("validate failed")
		}
		if trans := validation.Translate2Chinese(err); "[Name] "+validation.CharTypeGT3TransErr != trans {
			t.Errorf("validate failed.want:[%s] got:[%s]", validation.CharTypeGT3TransErr, trans)
			return
		}
		response.MkResponse(ctx, http.StatusOK, response.Success, "pong")
	})
	performRequest(r, "GET", "/ping?name=123A", nil)
}

func TestValidator_POST(t *testing.T) {
	r := testServer(t)
	r.POST("/ping", func(ctx *gin.Context) {
		c := validateCase{}
		err := ctx.ShouldBindJSON(&c)
		if err == nil {
			t.Errorf("validate failed")
		}
		if trans := validation.Translate2Chinese(err); "[Name] "+validation.CharTypeGT3TransErr != trans {
			t.Errorf("validate failed.want:[%s] got:[%s]", validation.CharTypeGT3TransErr, trans)
			return
		}
		response.MkResponse(ctx, http.StatusOK, response.Success, "pong")
	})
	b, _ := json.Marshal(validateCase{"123A"})
	performRequest(r, "POST", "/ping", bytes.NewReader(b), header{"Content-Type", binding.MIMEJSON})
}

type header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, body io.Reader, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	r.ServeHTTP(w, req)
	return w
}
