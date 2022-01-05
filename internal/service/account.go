package service

import (
	"context"
	"errors"
	"github.com/galaxy-future/BridgX/pkg/cloud/tencent"
	"strings"
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/errs"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/galaxy-future/BridgX/pkg"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/cloud/alibaba"
	"github.com/galaxy-future/BridgX/pkg/cloud/huawei"
	"github.com/galaxy-future/BridgX/pkg/encrypt"
)

// GetAccounts search accounts by condition.
func GetAccounts(provider, accountName, accountKey string, pageNum, pageSize int) ([]*model.Account, int64, error) {
	res, count, err := model.GetAccounts(provider, accountName, accountKey, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	DecryptAccounts(res)
	return res, count, nil
}

// decryptAccounts will decrypt account's `EncryptedAccountSecret` field.
// If strict param's value is true error will return.
func decryptAccounts(accounts []*model.Account, strict bool) error {
	if accounts == nil || len(accounts) == 0 {
		return nil
	}
	for _, account := range accounts {
		decrypted, err := DecryptAccount(encrypt.AesKeyPepper, account.Salt, account.AccountKey, account.EncryptedAccountSecret)
		if err != nil && strict == true {
			return err
		}
		account.AccountSecret = decrypted
	}
	return nil
}

// MustDecryptAccounts same as decryptAccounts(accounts, true).
func MustDecryptAccounts(accounts []*model.Account) error {
	return decryptAccounts(accounts, true)
}

// DecryptAccounts same as decryptAccounts(accounts, false).
func DecryptAccounts(accounts []*model.Account) {
	_ = decryptAccounts(accounts, false)
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
	err = MustDecryptAccounts(a)
	if err != nil {
		return nil, err
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
	err = MustDecryptAccounts([]*model.Account{a})
	if err != nil {
		return nil, err
	}
	account.Info = append(account.Info, types.KeyInfo{
		AK:       a.AccountKey,
		SK:       a.AccountSecret,
		Provider: a.Provider,
	})
	return &account, err
}

func CheckAccountValid(ak, sk, provider string) error {
	var err error
	var cli cloud.Provider
	switch provider {
	case cloud.AlibabaCloud:
		cli, err = alibaba.New(ak, sk, DefaultRegion)
	case cloud.HuaweiCloud:
		cli, err = huawei.New(ak, sk, DefaultRegionHuaWei)
	case cloud.TencentCloud:
		cli, err = tencent.New(ak, sk, DefaultRegionTencent)
	default:
		return errors.New("invalid provider")
	}

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
	err := createAccount(account)
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

func GetAksByOrgAk(ctx context.Context, orgId int64, ak string) ([]string, error) {
	return model.GetAksByOrgAk(ctx, orgId, ak)
}

func GetOrgKeysByAk(ctx context.Context, ak string) (*types.OrgKeys, error) {
	a, err := model.GetAccountByAk(ctx, ak)
	if err != nil {
		return nil, err
	}
	err = MustDecryptAccounts([]*model.Account{&a})
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

func EncryptAccount(pepper, salt, key, text string) (string, error) {
	encrypted, err := encrypt.AESEncrypt(key, wrapText(pepper, text, salt))
	if err != nil {
		return "", err
	}
	return encrypted, nil
}

func wrapText(pepper, text, salt string) string {
	return encrypt.ObfuscateText(pepper, text, salt)
}

func unWrapText(pepper, decryptedText, salt string) (string, error) {
	return encrypt.RestoreText(pepper, decryptedText, salt)
}

func DecryptAccount(pepper, salt, key, encrypted string) (string, error) {
	decrypted, err := encrypt.AESDecrypt(key, encrypted)
	if err != nil {
		return "", err
	}
	return unWrapText(pepper, decrypted, salt)
}

func generateSalt() (uid string, err error) {
	defer func() {
		if uid != "" {
			uid = strings.ReplaceAll(uid, "-", "")
		}
	}()
	uid, err = pkg.NewUUID()
	if err == nil && uid != "" {
		return uid, nil
	}
	return pkg.NewUUID4()
}

func createAccount(account *model.Account) error {
	salt, err := generateSalt()
	if err != nil {
		logs.Logger.Errorf("save account falied.because Because the salt generation failed.err: [%s]", err.Error())
		return errs.ErrSaveAccountFailed
	}
	encrypted, err := EncryptAccount(encrypt.AesKeyPepper, salt, account.AccountKey, account.AccountSecret)
	if err != nil {
		logs.Logger.Errorf("save account falied.because Because the account key secrect encryption failed.err: [%s]", err.Error())
		return errs.ErrSaveAccountFailed
	}

	account.Salt = salt
	account.EncryptedAccountSecret = encrypted
	return model.Save(account)
}

// GetAccountSecretByAccountKey get sk(decrypt) by ak
func GetAccountSecretByAccountKey(ctx context.Context, ak string) (string, error) {
	account, err := model.GetAccountByAk(ctx, ak)
	if err != nil {
		return "", err
	}
	sk, err := DecryptAccount(encrypt.AesKeyPepper, account.Salt, account.AccountKey, account.EncryptedAccountSecret)
	if err != nil {
		return "", err
	}
	return sk, nil
}
