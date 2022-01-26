package response

import "github.com/gin-gonic/gin"

func MkResponse(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

const (
	Success          = "success"
	ParamInvalid     = "param_invalid"
	TokenInvalid     = "token_invalid"
	PermissionDenied = "permission_denied"
	UserNotFound     = "user_not_found"
	UserExists       = "user_already_exists"
	TaskNotFound     = "task_not_found"
)
