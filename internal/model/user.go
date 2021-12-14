package model

import (
	"context"
	"fmt"

	"github.com/galaxy-future/BridgX/internal/cache"
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/logs"
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
	ret, _, err := GetUserThroughBigCache([]int64{uid}, cache.UserKeyMaker, func(ids []int64) ([]*User, error) {
		logs.Logger.Infof("get user:%v from db", uid)
		user := User{}
		err := Get(uid, &user)
		if err != nil {
			return nil, err
		}
		return []*User{&user}, nil
	})
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("uid:%v not found", uid)
	}
	return ret[0], nil
}

func GetUserThroughBigCache(ids []int64, keyMaker func(int64) string, delegate func(ids []int64) ([]*User, error)) ([]*User, []int64, error) {
	users := make([]*User, 0)
	needFetchIds, err := cache.GetFromBigCache(ids, &users, keyMaker)
	logs.Logger.Infof("Get user from local cache:%v", ids)
	if err != nil {
		needFetchIds = ids
	}
	if len(needFetchIds) > 0 && delegate != nil {
		logs.Logger.Infof("Get users:%v from delegate func", needFetchIds)
		rest, err := delegate(needFetchIds)
		if err != nil {
			return nil, ids, err
		}
		users = append(users, rest...)
		for _, user := range rest {
			logs.Logger.Infof("Set user to local cache uid:%v", keyMaker(user.GetId()))
			_ = cache.SetBigCache(user.GetId(), user, keyMaker)
		}
	}
	return users, needFetchIds, nil
}

func GetUsersByIDs(ctx context.Context, ids []int64) []User {
	users := make([]User, 0)
	err := clients.ReadDBCli.WithContext(ctx).
		Where("id in (?)", ids).Find(&users).
		Error
	if err != nil {
		logErr("get user from readDB", err)
		return nil
	}
	return users
}
