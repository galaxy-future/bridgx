package tests

import (
	"context"
	"testing"

	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/service"
)

func TestRecordOperationLog(t *testing.T) {
	type args struct {
		ctx   context.Context
		oplog service.OperationLog
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "diff",
			args: args{
				ctx: nil,
				oplog: service.OperationLog{
					Operation: "edit cluster",
					Operator:  2,
					Old: model.Account{
						Base:        model.Base{},
						AccountName: "1",
					},
					New: model.Account{
						Base:        model.Base{},
						AccountName: "2",
						UpdateBy:    "dasdad",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := service.RecordOperationLog(tt.args.ctx, tt.args.oplog); (err != nil) != tt.wantErr {
				t.Errorf("RecordOperationLog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
