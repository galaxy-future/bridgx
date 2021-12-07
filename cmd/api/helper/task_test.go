package helper

import (
	"testing"

	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/model"
)

func Test_getTaskInfoCountDiff(t *testing.T) {
	type args struct {
		task *model.Task
	}
	tests := []struct {
		name       string
		args       args
		wantBefore int
		wantExpect int
	}{
		{
			name: "expand instance count diff",
			args: args{
				task: &model.Task{
					TaskAction: constants.TaskActionExpand,
					TaskInfo:   "{\"count\":5,\"before_count\":10}",
				},
			},
			wantBefore: 10,
			wantExpect: 15,
		},
		{
			name: "shrink instance count diff",
			args: args{
				task: &model.Task{
					TaskAction: constants.TaskActionShrink,
					TaskInfo:   "{\"count\":6,\"before_count\":10}",
				},
			},
			wantBefore: 10,
			wantExpect: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBefore, gotExpect := getTaskInfoCountDiff(tt.args.task)
			if gotBefore != tt.wantBefore {
				t.Errorf("getTaskInfoCountDiff() gotBefore = %v, want %v", gotBefore, tt.wantBefore)
			}
			if gotExpect != tt.wantExpect {
				t.Errorf("getTaskInfoCountDiff() gotExpect = %v, want %v", gotExpect, tt.wantExpect)
			}
		})
	}
}
