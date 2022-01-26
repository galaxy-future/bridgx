package handler

import (
	"net/http"
	"strings"

	"github.com/galaxy-future/BridgX/internal/model"

	"github.com/galaxy-future/BridgX/internal/constants"

	"github.com/spf13/cast"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/cmd/api/middleware/authorization"
	"github.com/galaxy-future/BridgX/cmd/api/request"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/config"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	loginReq := request.LoginRequest{}
	if err := ctx.ShouldBind(&loginReq); err != nil {
		ctx.Abort()
		response.MkResponse(ctx, http.StatusBadRequest, "missing username/password", nil)
		return
	}
	user := service.Login(ctx, loginReq.Username, loginReq.Password)
	if user == nil {
		ctx.Abort()
		response.MkResponse(ctx, http.StatusBadRequest, "incorrect username/password", nil)
		return
	}
	if user.UserStatus == constants.UserStatusDisable {
		response.MkResponse(ctx, http.StatusBadRequest, "your account has been suspended", nil)
		return
	}
	userTokenFactory := authorization.CreateUserTokenFactory()
	userToken, err := userTokenFactory.GenerateToken(user.Id, user.Username, helper.ConvertToReadableStr(user.UserType), user.OrgId, config.GlobalConfig.JwtToken.JwtTokenCreatedExpires)
	if err == nil {
		response.MkResponse(ctx, http.StatusOK, response.Success, userToken)
		return
	}
	response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	return
}

func RefreshToken(ctx *gin.Context) {
	headerParams := authorization.HeaderParams{}
	if err := ctx.ShouldBindHeader(&headerParams); err != nil {
		ctx.Abort()
		response.MkResponse(ctx, http.StatusBadRequest, "missing jwt token", nil)
		return
	}
	token := strings.Split(headerParams.Authorization, " ")
	if len(token) == 2 && len(token[1]) >= 20 {
		isTokenValid := authorization.CreateUserTokenFactory().IsValid(token[1])
		if isTokenValid {
			if newToken, ok := authorization.CreateUserTokenFactory().RefreshToken(token[1]); ok {
				response.MkResponse(ctx, http.StatusOK, response.Success, newToken)
				return
			}
		} else {
			response.MkResponse(ctx, http.StatusUnauthorized, "bad token refresh", nil)
			ctx.Abort()
			return
		}
	} else {
		response.MkResponse(ctx, http.StatusUnauthorized, "token base info error", nil)
		ctx.Abort()
	}
}

func GetUserInfo(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}
	userInDB, err := service.GetUserById(ctx, user.UserId)
	if err != nil || userInDB == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.UserNotFound, nil)
		return
	}
	res := &response.UserInfo{
		UserId:   user.UserId,
		Username: userInDB.Username,
		OrgId:    userInDB.OrgId,
		UserType: helper.ConvertToReadableStr(userInDB.UserType),
		CreateAt: userInDB.CreateAt.Format("2006-01-02 15:04:05"),
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, res)
	return
}

func ListUsers(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}

	pn, ps := getPager(ctx)

	users, total, err := service.GetUserList(ctx, user.OrgId, pn, ps)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	pager := response.Pager{
		PageNumber: pn,
		PageSize:   ps,
		Total:      int(total),
	}
	resp := &response.ListUsersResponse{
		UserList: helper.ConvertToUserThumbList(users),
		Pager:    pager,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func CreateUser(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}
	req := request.CreateUserRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	dbUser := model.GetUserByName(ctx, req.UserName)
	if dbUser != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.UserExists, nil)
		return
	}
	if user.UserType == constants.UserTypeCommonUserStr || (user.Name != constants.UserNameRoot && req.UserType == constants.UserTypeAdminStr) {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}
	err := service.CreateUser(ctx, user.OrgId, req.UserName, req.Password, user.Name, helper.ConvertUserTypeToInt(req.UserType))
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

func EnableUser(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}

	req := request.UserStatusRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}

	if user.UserType == constants.UserTypeCommonUserStr {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}

	if user.Name != constants.UserNameRoot && user.UserType == constants.UserTypeAdminStr {
		exist, err := service.ExistAdmin(ctx, req.UserNames)
		if err != nil {
			response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		if exist {
			response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
			return
		}

	}

	err := service.UpdateUserStatus(ctx, req.UserNames, req.Action)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

func ModifyAdminPassword(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}

	req := request.ModifyAdminPasswordRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}

	err := service.ModifyAdminPassword(ctx, user.UserId, user.Name, req.OldPassword, req.NewPassword)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

func ModifyUsername(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}

	req := request.ModifyUsernameRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	if req.NewUsername == "" {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	if cast.ToInt64(req.UserId) != user.UserId {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}

	err := service.ModifyUsername(ctx, user.UserId, req.NewUsername)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

func ModifyUsertype(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}

	req := request.ModifyUserTypeRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	if user.Name != constants.UserNameRoot {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}

	err := service.ModifyUsertype(ctx, req.UserIds, helper.ConvertUserTypeToInt(req.UserType))
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}
