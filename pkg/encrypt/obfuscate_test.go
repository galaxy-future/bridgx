package encrypt

import (
	"testing"
)

func TestObfuscateRestoreText(t *testing.T) {
	type args struct {
		pepper string
		text   string
		salt   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "obfuscate then restore",
			args: args{
				pepper: "bridgx",
				text:   "asdfghjklqwertyuiop",
				salt:   "9ef8222b-123c-4c14-beea-49ea7b8da03d",
			},
			want: "262796163746667686a6b6c619356668323232326d213233336d243361343467687d226565616d2439356167326834616033346777562747975796f607",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ObfuscateText(tt.args.pepper, tt.args.text, tt.args.salt)
			if got != tt.want {
				t.Errorf("ObfuscateText() = %v, want %v", got, tt.want)
			}
			text, err := RestoreText(tt.args.pepper, got, tt.args.salt)
			if err != nil {
				t.Errorf("RestoreText failed %s", err.Error())
			}
			if text != tt.args.text {
				t.Errorf("RestoreText() = %v, want %v", tt.args.text, text)
			}
		})
	}
}

func Test_swapOddEven(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				t: "",
			},
			want: "",
		},
		{
			name: "even",
			args: args{
				t: "abcd",
			},
			want: "badc",
		},
		{
			name: "odd",
			args: args{
				t: "abcde",
			},
			want: "badce",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := swapOddEven(tt.args.t); got != tt.want {
				t.Errorf("swapOddEven() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObfuscateText(t *testing.T) {
	type args struct {
		pepper string
		text   string
		salt   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{},
			want: "",
		},
		{
			name: "empty pepper",
			args: args{
				pepper: "",
				text:   "1",
				salt:   "1",
			},
			want: "1133",
		},
		{
			name: "empty text",
			args: args{
				pepper: "1",
				text:   "",
				salt:   "1",
			},
			want: "1133",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ObfuscateText(tt.args.pepper, tt.args.text, tt.args.salt); got != tt.want {
				t.Errorf("ObfuscateText() = %v, want %v", got, tt.want)
			}
		})
	}
}
