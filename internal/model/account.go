package model

import (
	"context"
	"fmt"

	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/encrypt"

	"github.com/galaxy-future/BridgX/internal/clients"
	"gorm.io/gorm"
)

// Account cloud provider account info
type Account struct {
	Base
	AccountName   string `json:"account_name"`
	AccountKey    string `json:"account_key"`
	AccountSecret string `json:"account_secret"`
	Provider      string `json:"provider"`
	OrgId         int64  `json:"org_id"`
	CreateBy      string `json:"create_by"`
	UpdateBy      string `json:"update_by"`
	DeletedAt     gorm.DeletedAt
}

// TableName table name in DB
func (a Account) TableName() string {
	return "account"
}

//AfterFind decrypt account secret
func (a *Account) AfterFind(tx *gorm.DB) (err error) {
	if a == nil {
		return nil
	}
	if a.AccountKey != "" && a.AccountSecret != "" {
		res, err := encrypt.AESDecrypt(a.AccountKey+encrypt.AesKeySalt, a.AccountSecret)
		if err != nil {
			logs.Logger.Errorf("decrypt sk failed.err: %s", err.Error())
			return err
		}
		a.AccountSecret = res
	}
	return nil
}

//BeforeSave encrypt account secret before insert DB
func (a *Account) BeforeSave(tx *gorm.DB) (err error) {
	if a == nil {
		return nil
	}
	if a.AccountKey != "" && a.AccountSecret != "" {
		res, err := encrypt.AESEncrypt(a.AccountKey+encrypt.AesKeySalt, a.AccountSecret)
		if err != nil {
			logs.Logger.Errorf("encrypt sk failed.err: %s", err.Error())
			return err
		}
		a.AccountSecret = res
	}
	return nil
}

//GetAccounts search accounts by condition
func GetAccounts(provider, accountName, accountKey string, pageNum, pageSize int) ([]Account, int64, error) {
	res := make([]Account, 0)

	sql := clients.ReadDBCli.Where(map[string]interface{}{})
	if accountName != "" {
		sql.Where("account_name LIKE ?", fmt.Sprintf("%%%v%%", accountName))
	}
	if provider != "" {
		sql.Where("provider = ?", provider)
	}
	if accountKey != "" {
		sql.Where("account_key = ?", accountKey)
	}
	err := sql.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&res).Error
	if err != nil {
		return res, 0, err
	}
	var cnt int64
	err = sql.Offset(-1).Limit(-1).Count(&cnt).Error
	if err != nil {
		return res, 0, err
	}
	return res, cnt, err
}

//GetAccountSecretByAccountKey get sk(decrypt) by ak
func GetAccountSecretByAccountKey(ak string) string {
	var ac Account
	if err := clients.ReadDBCli.Where("account_key = ?", ak).First(&ac).Error; err != nil {
		logErr("GetAccountSecretByAccountKey from read db", err)
		return ""
	}
	return ac.AccountSecret
}

//GetAccountsByOrgId get accounts belongs to specify orgId
func GetAccountsByOrgId(orgId int64) (accounts []Account, err error) {
	if err := clients.ReadDBCli.Where("org_id = ?", orgId).Find(&accounts).Error; err != nil {
		logErr("GetAccountSecretByAccountKey from read db", err)
		return nil, err
	}
	return accounts, nil
}

//GetDefaultAccountByProvider return default accounts by provider
func GetDefaultAccountByProvider(provider string) (account Account, err error) {
	if err := clients.ReadDBCli.Where("provider = ?", provider).First(&account).Error; err != nil {
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

// GetAccountsByAk get first account by ak
func GetAccountsByAk(ctx context.Context, ak string) (a Account, err error) {
	err = clients.ReadDBCli.WithContext(ctx).
		Where("account_key = ?", ak).
		First(&a).Error
	if err != nil {
		return Account{}, err
	}
	return a, nil
}
