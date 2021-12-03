package service

import (
	"context"
	"errors"
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/galaxy-future/BridgX/pkg/cloud/alibaba"
)

func GetAccounts(provider, accountName, accountKey string, pageNum, pageSize int) ([]model.Account, int64, error) {
	res, total, err := model.GetAccounts(provider, accountName, accountKey, pageNum, pageSize)
	return res, total, err
}

//GetAccount query account info by provider and accountKey
func GetAccount(provider, accountKey string) (*model.Account, error) {
	var ret model.Account
	err := model.QueryFirst(map[string]interface{}{"provider": provider, "account_key": accountKey}, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetAccountsByOrgId(orgId int64) (*types.OrgKeys, error) {
	a, err := model.GetAccountsByOrgId(orgId)
	if err != nil {
		return nil, err
	}
	account := types.OrgKeys{
		OrgId: orgId,
	}
	for _, info := range a {
		account.Info = append(account.Info, types.KeyInfo{
			AK:       info.AccountKey,
			SK:       info.AccountSecret,
			Provider: info.Provider,
		})
	}
	return &account, err
}

func GetDefaultAccount(provider string) (*types.OrgKeys, error) {
	a, err := model.GetDefaultAccountByProvider(provider)
	if err != nil {
		return nil, err
	}
	account := types.OrgKeys{
		OrgId: a.OrgId,
	}
	account.Info = append(account.Info, types.KeyInfo{
		AK:       a.AccountKey,
		SK:       a.AccountSecret,
		Provider: a.Provider,
	})
	return &account, err
}

func CheckAccountValid(ak, sk string) error {
	cli, err := alibaba.New(ak, sk, DefaultRegion)
	if err != nil {
		return err
	}
	_, err = cli.GetRegions()
	return err
}

func CreateCloudAccount(ctx context.Context, accountName, provider, ak, sk string, orgId int64, username string) error {
	account := &model.Account{
		AccountName:   accountName,
		AccountKey:    ak,
		AccountSecret: sk,
		Provider:      provider,
		OrgId:         orgId,
		CreateBy:      username,
		UpdateBy:      username,
	}
	now := time.Now()
	account.CreateAt = &now
	account.UpdateAt = &now
	err := model.Create(account)
	if err != nil {
		return err
	}
	H.SubmitTask(&SimpleTask{
		ProviderName: provider,
		AccountKey:   ak,
		TargetType:   TargetTypeAccount,
		Retry:        3,
	})
	H.SubmitTask(&SimpleTask{
		ProviderName: provider,
		AccountKey:   ak,
		TargetType:   TargetTypeInstanceType,
		Retry:        5,
	})
	return nil
}

func EditCloudAccount(ctx context.Context, id int64, accountName, provider, username string) error {
	account := model.Account{}
	err := model.Get(id, &account)
	if err != nil {
		return err
	}
	if accountName != "" {
		account.AccountName = accountName
	}
	if provider != "" {
		account.Provider = provider
	}
	account.UpdateBy = username
	now := time.Now()
	account.UpdateAt = &now
	return model.Save(&account)
}

func DeleteCloudAccount(ctx context.Context, ids []int64, orgId int64) error {
	accounts := make([]model.Account, 0)
	if len(ids) == 0 {
		return nil
	}
	err := model.Gets(ids, &accounts)
	if err != nil {
		return err
	}
	if len(accounts) == 0 {
		return nil
	}
	for _, account := range accounts {
		if account.OrgId != orgId {
			return errors.New("delete permission denied")
		}
	}
	err = clients.WriteDBCli.WithContext(ctx).Delete(&accounts).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAksByOrgId(orgId int64) ([]string, error) {
	accounts, err := model.GetAccountsByOrgId(orgId)
	if err != nil {
		return nil, err
	}
	aks := make([]string, 0, len(accounts))
	for _, a := range accounts {
		aks = append(aks, a.AccountKey)
	}
	return aks, nil
}

func GetAksByOrgAkProvider(ctx context.Context, orgId int64, ak, provider string) ([]string, error) {
	return model.GetAksByOrgAkProvider(ctx, orgId, ak, provider)
}

func GetOrgKeysByAk(ctx context.Context, ak string) (*types.OrgKeys, error) {
	a, err := model.GetAccountsByAk(ctx, ak)
	if err != nil {
		return nil, err
	}
	return &types.OrgKeys{
		OrgId: 0,
		Info: []types.KeyInfo{{
			AK:       a.AccountKey,
			SK:       a.AccountSecret,
			Provider: a.Provider,
		}},
	}, nil
}
