package cmp

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/galaxy-future/BridgX/pkg/utils"
	"github.com/google/go-cmp/cmp"
)

type user struct {
	Name     string            `json:"name" diff:"-"`
	Age      int               `json:"age" diff:"age"`
	Birthday *time.Time        `json:"birthday" diff:"birthday"`
	M        map[string]string `json:"m" diff:"m"`
}

func Test_compare(t *testing.T) {
	type args struct {
		t  reflect.Type
		vx reflect.Value
		vy reflect.Value
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "different kind",
			args: args{
				t:  nil,
				vx: reflect.ValueOf(0),
				vy: reflect.ValueOf(""),
			},
			want: false,
		},
		{
			name: "same int",
			args: args{
				t:  nil,
				vx: reflect.ValueOf(0),
				vy: reflect.ValueOf(0),
			},
			want: true,
		},
		{
			name: "same float",
			args: args{
				t:  nil,
				vx: reflect.ValueOf(0.1),
				vy: reflect.ValueOf(0.10),
			},
			want: true,
		},
		{
			name: "same string",
			args: args{
				t:  nil,
				vx: reflect.ValueOf("1"),
				vy: reflect.ValueOf("1"),
			},
			want: true,
		},
		{
			name: "same slice",
			args: args{
				t:  nil,
				vx: reflect.ValueOf([]int{1, 2}),
				vy: reflect.ValueOf([]int{1, 2}),
			},
			want: true,
		},
		{
			name: "same map",
			args: args{
				t:  nil,
				vx: reflect.ValueOf(map[int]int{1: 2}),
				vy: reflect.ValueOf(map[int]int{1: 2}),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compare(tt.args.vx, tt.args.vy); got != tt.want {
				t.Errorf("compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiff(t *testing.T) {
	t1 := time.Now()
	t2 := time.Now().Add(time.Minute)
	type args struct {
		old interface{}
		new interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantDiff []map[string]string
		wantErr  bool
	}{
		{
			name: "same empty struct",
			args: args{
				old: user{},
				new: user{},
			},
			wantDiff: []map[string]string{},
			wantErr:  false,
		},
		{
			name: "one different field",
			args: args{
				old: user{
					Name:     "张三",
					Age:      21,
					Birthday: &t1,
				},
				new: user{
					Name:     "张三",
					Age:      20,
					Birthday: &t1,
				},
			},
			wantDiff: []map[string]string{
				{"operation": "edit", "target": "age", "old": "21", "new": "20"},
			},
			wantErr: false,
		},
		{
			name: "multi different field",
			args: args{
				old: user{
					Name:     "张三",
					Age:      21,
					Birthday: &t1,
				},
				new: user{
					Name:     "张三",
					Age:      20,
					Birthday: &t2,
				},
			},
			wantDiff: []map[string]string{
				{"operation": "edit", "target": "age", "old": "21", "new": "20"},
				{"target": "birthday", "old": utils.FormatTime(t1), "new": utils.FormatTime(t2), "operation": "edit"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := Diff(tt.args.old, tt.args.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("Diff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			display, err := gotRes.Beautiful()
			if err != nil {
				t.Errorf("Beautiful() error = %v", err)
				return
			}
			if !reflect.DeepEqual(display, tt.wantDiff) {
				fmt.Println(cmp.Diff(display, tt.wantDiff))
				t.Errorf("Diff() display = %v, want %v", display, tt.wantDiff)
			}
		})
	}
}
