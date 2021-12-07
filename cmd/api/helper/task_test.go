package helper

import (
	"testing"

	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/model"
)

func Test_getTaskInfoCountDiff(t *testing.T) {
	success := 1
	type args struct {
		task *model.Task
	}
	tests := []struct {
		name       string
		args       args
		wantBefore int
		wantExpect int
		wantAfter  int
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
			wantAfter:  11,
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
			wantAfter:  9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBefore, gotAfter, gotExpect := getTaskInfoCountDiff(tt.args.task, success)
			if gotBefore != tt.wantBefore {
				t.Errorf("getTaskInfoCountDiff() gotBefore = %v, want %v", gotBefore, tt.wantBefore)
			}
			if gotExpect != tt.wantExpect {
				t.Errorf("getTaskInfoCountDiff() gotExpect = %v, want %v", gotExpect, tt.wantExpect)
			}
			if gotAfter != tt.wantAfter {
				t.Errorf("getTaskInfoCountDiff() gotAfter = %v, want %v", gotAfter, tt.wantAfter)
			}
		})
	}
}
