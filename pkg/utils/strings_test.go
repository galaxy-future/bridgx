package utils

import (
	"reflect"
	"testing"
)

func TestStringSliceSplit(t *testing.T) {
	type args struct {
		slice     []string
		singleLen int64
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "Test1",
			args: args{
				slice:     []string{"a", "b", "c", "d", "e"},
				singleLen: 2,
			},
			want: [][]string{
				{"a", "b"},
				{"c", "d"},
				{"e"},
			},
		},
		{
			name: "Test2",
			args: args{
				slice:     []string{"a", "b", "c", "d", "e"},
				singleLen: 5,
			},
			want: [][]string{
				{"a", "b", "c", "d", "e"},
			},
		},
		{
			name: "Test3",
			args: args{
				slice:     []string{"a", "b", "c", "d", "e"},
				singleLen: 6,
			},
			want: [][]string{
				{"a", "b", "c", "d", "e"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceSplit(tt.args.slice, tt.args.singleLen); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSliceSplit() = %v, want %v", got, tt.want)
			}
		})
	}
}
