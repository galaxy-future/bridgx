package authorization

import (
	"net/http"
	"strings"

	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"

	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/config"
	"github.com/gin-gonic/gin"
)

type HeaderParams struct {
	Authorization string `header:"Authorization" binding:"required,min=20"`
}

func CheckTokenAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headerParams := HeaderParams{}
		if err := ctx.ShouldBindHeader(&headerParams); err != nil {
			ctx.Abort()
			response.MkResponse(ctx, http.StatusBadRequest, "missing jwt token", nil)
			return
		}
		token := strings.Split(headerParams.Authorization, " ")
		if len(token) == 2 && len(token[1]) >= 20 {
			isTokenValid := CreateUserTokenFactory().IsValid(token[1])
			if isTokenValid {
				ctx.Set(gf_cluster.HeaderTokenName, token[1])
				if customToken, err := CreateUserTokenFactory().ParseToken(token[1]); err == nil {
					key := config.GlobalConfig.JwtToken.BindContextKeyName
					ctx.Set(key, customToken)
				}
				ctx.Next()
			} else {
				response.MkResponse(ctx, http.StatusUnauthorized, "token auth fail", nil)
				ctx.Abort()
				return
			}
		} else {
			response.MkResponse(ctx, http.StatusUnauthorized, "token base info error", nil)
			ctx.Abort()
		}
	}
}

func RefreshTokenConditionCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headerParams := HeaderParams{}
		if err := ctx.ShouldBindHeader(&headerParams); err != nil {
			response.MkResponse(ctx, http.StatusBadRequest, "missing jwt token", nil)
			ctx.Abort()
			return
		}
		token := strings.Split(headerParams.Authorization, " ")
		if len(token) == 2 && len(token[1]) >= 20 {
			// 判断token是否满足刷新条件
			if CreateUserTokenFactory().TokenIsMeetRefreshCondition(token[1]) {
				ctx.Next()
			} else {
				response.MkResponse(ctx, http.StatusUnauthorized, "token refresh fail", nil)
				ctx.Abort()
				return
			}
		} else {
			response.MkResponse(ctx, http.StatusUnauthorized, "token refresh fail", nil)
			ctx.Abort()
			return
		}
	}
}
