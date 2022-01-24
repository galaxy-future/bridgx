package helper

import (
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/spf13/cast"
)

func ConvertToUserThumbList(users []model.User) []response.UserThumb {
	res := make([]response.UserThumb, 0)
	for _, user := range users {
		c := response.UserThumb{
			UserId:     cast.ToString(user.Id),
			UserName:   user.Username,
			CreateAt:   user.CreateAt.Format("2006-01-02 15:04:05"),
			CreateBy:   user.CreateBy,
			UserStatus: user.UserStatus,
			UserType:   ConvertToReadableStr(user.UserType),
		}
		res = append(res, c)
	}
	return res
}

func ConvertToReadableStr(userType int8) string {
	if userType == constants.UserTypeAdmin {
		return "ADMIN"
	}
	if userType == constants.UserTypeCommonUser {
		return "COMMON"
	}
	return "SYSTEM"
}

func ConvertUserTypeToInt(userType string) int8 {
	if userType == "ADMIN" {
		return constants.UserTypeAdmin
	}
	return constants.UserTypeCommonUser
}
