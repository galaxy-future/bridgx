package tests

import (
	"testing"

	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGetAccounts(t *testing.T) {
	accounts, total, err := service.GetAccounts("AlibabaCloud", "TES", "", 1, 10)
	assert.Nil(t, err)
	assert.Len(t, accounts, 2)
	assert.EqualValues(t, total, 2)
}

func TestGetAccount(t *testing.T) {
	account, err := service.GetAccount("AlibabaCloud", "TES")
	assert.Nil(t, err)
	assert.NotNil(t, account)

	account, err = service.GetAccount("AlibabaCloud", "xxxx")
	assert.NotNil(t, err)
	assert.Nil(t, account)
}
