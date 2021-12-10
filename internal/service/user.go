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

func GetUserList(ctx context.Context, orgId int64, pageNum, pageSize int) (ret []model.User, total int64, err error) {
	queryMap := map[string]interface{}{"org_id": orgId, "user_type": []int{constants.UserTypeCommonUser}}

	total, err = model.Query(queryMap, pageNum, pageSize, &ret, "id DESC", true)
	if err != nil {
		return ret, 0, err
	}
	return ret, total, nil
}

func CreateUser(ctx context.Context, orgId int64, username, password, createBy string) error {
	user := &model.User{
		Username:   username,
		Password:   utils.Base64Md5(password),
		OrgId:      orgId,
		UserStatus: constants.UserStatusEnable,
		UserType:   constants.UserTypeCommonUser,
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
