package encrypt

import (
	"testing"
)

func TestAESDecrypt(t *testing.T) {
	type args struct {
		key  string
		ct16 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"decrypt",
			args{
				key:  "bridgx",
				ct16: "ec2d948a21ecd5868057bada5b315447",
			},
			"xxx",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AESDecrypt(tt.args.key, tt.args.ct16)
			if (err != nil) != tt.wantErr {
				t.Errorf("AESDecrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AESDecrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAESEncrypt(t *testing.T) {
	type args struct {
		key       string
		plaintext string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"encrypt",
			args{
				key:       "bridgx",
				plaintext: "xxx",
			},
			"ec2d948a21ecd5868057bada5b315447",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AESEncrypt(tt.args.key, tt.args.plaintext)
			if (err != nil) != tt.wantErr {
				t.Errorf("AESEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AESEncrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
