package utils

import "testing"

func TestNewErrf(t *testing.T) {
	type args struct {
		format string
		a      []interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantErrMsg string
	}{
		{
			name: "format errors",
			args: args{
				format: "%derr%s%d",
				a: []interface{}{
					1, "ssss", 2,
				},
			},
			wantErrMsg: "1errssss2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewErrf(tt.args.format, tt.args.a...); (err != nil) && err.Error() != tt.wantErrMsg {
				t.Errorf("NewErrf() error = %v, wantErr %v", err, tt.wantErrMsg)
			}
		})
	}
}
