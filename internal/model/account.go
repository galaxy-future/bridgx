package model

import (
	"context"
	"fmt"

	"github.com/galaxy-future/BridgX/internal/clients"
	"gorm.io/gorm"
)

// Account cloud provider account info
type Account struct {
	Base
	AccountName            string `json:"account_name" diff:"account_name"`
	AccountKey             string `json:"account_key" diff:"account_key"`
	EncryptedAccountSecret string `json:"encrypted_account_secret" diff:"encrypted_account_secret"`
	Salt                   string `json:"salt" diff:"salt"`
	Provider               string `json:"provider" diff:"provider"`
	OrgId                  int64  `json:"org_id" diff:"org_id"`
	CreateBy               string `json:"create_by" diff:"create_by"`
	UpdateBy               string `json:"update_by" diff:"update_by"`
	DeletedAt              gorm.DeletedAt

	// the value of this field will not be empty only after decryption function called.
	AccountSecret string `json:"account_secret" gorm:"-"`
}

// TableName table name in DB
func (a Account) TableName() string {
	return "account"
}

//GetAccounts search accounts by condition
func GetAccounts(provider, accountName, accountKey string, pageNum, pageSize int) ([]*Account, int64, error) {
	res := make([]*Account, 0)
	query := clients.ReadDBCli.Table(Account{}.TableName())
	if accountName != "" {
		query.Where("account_name LIKE ?", fmt.Sprintf("%%%v%%", accountName))
	}
	if provider != "" {
		query.Where("provider = ?", provider)
	}
	if accountKey != "" {
		query.Where("account_key = ?", accountKey)
	}
	count, err := QueryWhere(query, pageNum, pageSize, &res, "id Desc", true)
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}

//GetAccountsByOrgId get accounts belongs to specify orgId
func GetAccountsByOrgId(orgId int64) (accounts []*Account, err error) {
	if err := clients.ReadDBCli.Where("org_id = ?", orgId).Find(&accounts).Error; err != nil {
		logErr("GetAccountSecretByAccountKey from read db", err)
		return nil, err
	}
	return accounts, nil
}

//GetDefaultAccountByProvider return default accounts by provider
func GetDefaultAccountByProvider(provider string) (account *Account, err error) {
	account = &Account{}
	if err := clients.ReadDBCli.Where("provider = ?", provider).First(account).Error; err != nil {
		logErr("GetAccountSecretByAccountKey from read db", err)
		return account, err
	}
	return account, nil
}

//GetAksByOrgAkProvider get aks by ak and provider
func GetAksByOrgAkProvider(ctx context.Context, orgId int64, ak, provider string) ([]string, error) {
	aks := make([]string, 0)
	query := clients.ReadDBCli.WithContext(ctx).
		Table(Account{}.TableName()).
		Select("account_key").
		Where("org_id = ?", orgId)
	if ak != "" {
		query = query.Where("account_key = ?", ak)
	}
	if provider != "" {
		query = query.Where("provider = ?", provider)
	}

	if err := query.Find(&aks).Error; err != nil {
		logErr("GetAksByOrgAkProvider from read db", err)
		return nil, err
	}
	return aks, nil
}

// GetAccountByAk get first account by ak
func GetAccountByAk(ctx context.Context, ak string) (a Account, err error) {
	err = clients.ReadDBCli.WithContext(ctx).
		Where("account_key = ?", ak).
		First(&a).Error
	if err != nil {
		return Account{}, err
	}
	return a, nil
}

func GetAllProvider(ctx context.Context) (provider []string, err error) {
	err = clients.ReadDBCli.WithContext(ctx).
		Model(Account{}).
		Select("provider").
		Group("provider").
		Find(&provider).Error
	if err != nil {
		return nil, err
	}
	return provider, nil
}
