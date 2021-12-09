package model

import (
	"context"

	"github.com/galaxy-future/BridgX/internal/clients"
)

type User struct {
	Base
	Username   string `json:"username"`
	Password   string `json:"password"`
	UserType   int8   `json:"user_type"`
	UserStatus string `json:"user_status"`
	OrgId      int64  `json:"org_id"`
	CreateBy   string `json:"create_by"`
}

func (u *User) TableName() string {
	return "user"
}

func GetUserByName(ctx context.Context, username string) *User {
	user := User{}
	err := clients.ReadDBCli.WithContext(ctx).Where(&User{Username: username}).Find(&user).Error
	if err != nil {
		logErr("get user from readDB", err)
		return nil
	}
	return &user
}

func UpdateUserStatus(ctx context.Context, model interface{}, usernames []string, updates map[string]interface{}) error {
	if err := clients.WriteDBCli.WithContext(ctx).Model(model).Where("username IN (?)", usernames).Updates(updates).Error; err != nil {
		logErr("update data list to write db", err)
		return err
	}
	return nil
}

func GetUserById(ctx context.Context, uid int64) (*User, error) {
	user := User{}
	err := Get(uid, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
