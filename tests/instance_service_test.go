package tests

import (
	"context"
	"testing"
	"time"

	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/galaxy-future/BridgX/pkg/cloud/alibaba"
)

func TestSyncInstanceTypes(t *testing.T) {
	type args struct {
		ctx      context.Context
		provider string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx:      context.Background(),
				provider: alibaba.CloudName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := service.SyncInstanceTypes(tt.args.ctx, tt.args.provider); (err != nil) != tt.wantErr {
				t.Errorf("SyncInstanceTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetInstanceUsageStatistics(t *testing.T) {
	day, _ := time.Parse("2006-01-02", "2021-11-12")
	res, total, err := service.GetInstanceUsageStatistics(context.Background(), "gf.bridgx.online", day, 0, 1, 10)
	t.Log(res)
	t.Log(total)
	t.Log(err)
}
