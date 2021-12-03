package helper

import (
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/pkg/encrypt"
	"github.com/spf13/cast"
)

// ConvertToCloudAccountList convert to account list display format
func ConvertToCloudAccountList(accounts []model.Account) []response.CloudAccount {
	res := make([]response.CloudAccount, 0)
	if len(accounts) == 0 {
		return res
	}
	for _, account := range accounts {
		ca := response.CloudAccount{
			Id:          cast.ToString(account.Id),
			AccountName: account.AccountName,
			AccountKey:  account.AccountKey,
			Provider:    account.Provider,
			CreateAt:    account.CreateAt.String(),
			CreateBy:    account.CreateBy,
		}
		res = append(res, ca)
	}
	return res
}

// ConvertToEncryptAccountInfo convert account_secret to encrypt account_secret
func ConvertToEncryptAccountInfo(account *model.Account) (*response.EncryptCloudAccountInfo, error) {
	accountSecretEncrypt, err := encrypt.AESEncrypt(account.AccountKey+encrypt.AesKeySalt, account.AccountSecret)
	if err != nil {
		return nil, err
	}
	return &response.EncryptCloudAccountInfo{
		AccountName:          account.AccountName,
		AccountKey:           account.AccountKey,
		AccountSecretEncrypt: accountSecretEncrypt,
		Provider:             account.Provider,
	}, nil
}
