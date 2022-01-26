package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/galaxy-future/BridgX/internal/constants"

	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/pkg/utils"
)

func Login(ctx context.Context, username, password string) *model.User {
	user := model.GetUserByName(ctx, username)
	if user == nil {
		return nil
	}
	if user.Password == utils.Base64Md5(password) {
		return user
	}
	return nil
}

func GetUserList(ctx context.Context, orgId int64, pageNum, pageSize int) ([]model.User, int64, error) {
	ret, total, err := model.GetUserList(ctx, orgId, pageNum, pageSize)
	if err != nil {
		return ret, 0, err
	}
	return ret, total, nil
}

func CreateUser(ctx context.Context, orgId int64, username, password, createBy string, userType int8) error {
	user := &model.User{
		Username:   username,
		Password:   utils.Base64Md5(password),
		OrgId:      orgId,
		UserStatus: constants.UserStatusEnable,
		UserType:   userType,
		CreateBy:   createBy,
	}
	now := time.Now()
	user.CreateAt = &now
	user.UpdateAt = &now
	return model.Create(user)
}

func UpdateUserStatus(ctx context.Context, usernames []string, status string) error {
	err := model.UpdateUserStatus(ctx, model.User{}, usernames, map[string]interface{}{"user_status": status, "update_at": time.Now()})
	if err != nil {
		return fmt.Errorf("can not update user stauts : %w", err)
	}
	return nil
}

func ExistAdmin(ctx context.Context, usernames []string) (bool, error) {
	users, err := model.GetUsersByUsernamesAndUserType(ctx, usernames, constants.UserTypeAdmin)
	if err != nil {
		return false, err
	}
	return len(users) > 0, nil

}

func ModifyAdminPassword(ctx context.Context, userId int64, userName, oldPassword, newPassword string) error {
	user := model.User{}
	err := model.Get(userId, &user)
	if err != nil {
		return err
	}

	if user.Password != utils.Base64Md5(oldPassword) {
		return fmt.Errorf("user old password does not match : %s", userName)
	}

	if newPassword != "" {
		user.Password = utils.Base64Md5(newPassword)
	}

	now := time.Now()
	user.UpdateAt = &now
	return model.Save(user)
}

func CountUser(ctx context.Context, orgId int64) (int64, error) {
	var ret []model.User
	return model.Count(map[string]interface{}{"org_id": orgId}, &ret)
}

func GetUserById(ctx context.Context, uid int64) (*model.User, error) {
	return model.GetUserById(ctx, uid)
}

func ModifyUsername(ctx context.Context, uid int64, newUsername string) error {
	user, err := GetUserById(ctx, uid)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	now := time.Now()
	user.Username = newUsername
	user.UpdateAt = &now
	return model.Save(user)
}

func ModifyUsertype(ctx context.Context, userIds []int64, userType int8) error {
	err := model.UpdateUserType(ctx, userIds, map[string]interface{}{"user_type": userType, "update_at": time.Now()})
	if err != nil {
		return err
	}
	return nil

}

func UserMapByIDs(ctx context.Context, ids []int64) map[int64]string {
	userMap := make(map[int64]string)
	users := model.GetUsersByIDs(ctx, ids)
	for _, user := range users {
		userMap[user.Id] = user.Username
	}
	return userMap
}
