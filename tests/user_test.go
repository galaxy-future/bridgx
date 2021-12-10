package tests

import (
	"context"
	"testing"

	"github.com/galaxy-future/BridgX/internal/model"
)

func TestGetUserById(t *testing.T) {
	for i := 0; i < 100; i++ {
		u, e := model.GetUserById(context.Background(), 1)
		if e != nil {
			t.Logf("Error:%v", e)
		}
		t.Logf("user:%v", *u)
	}
}
