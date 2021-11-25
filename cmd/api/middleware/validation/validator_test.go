package validation

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
)

const (
	sevenNumbers     = "1234567"
	thirtyOneNumbers = "1234567890123456789012345678901"
	tenNumbers       = "1234567890"
	tenUpperLetters  = "aaaaaaaaaa"
	tenLowerLetters  = "AAAAAAAAAA"
	tenSpecialChar   = "??????????"
)

func testValidator() (*validator.Validate, error) {
	v := validator.New()
	RegisterCustomValidators()
	err := RegisterValidators(v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func Test_validateAllCharacter(t *testing.T) {
	type args struct {
		Password string
	}
	errMsg := fmt.Sprintf("[Password] %s", CharTypeGT3TransErr)
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "验证长度小于8",
			args: args{
				Password: sevenNumbers,
			},
			want: "Password长度必须至少为8个字符",
		},
		{
			name: "验证长度大于30",
			args: args{
				Password: thirtyOneNumbers,
			},
			want: "Password长度不能超过30个字符",
		},
		{
			name: "验证纯数字",
			args: args{
				Password: tenNumbers,
			},
			want: errMsg,
		},
		{
			name: "验证纯大写字母",
			args: args{
				Password: tenUpperLetters,
			},
			want: errMsg,
		},
		{
			name: "验证纯小写字母",
			args: args{
				Password: tenLowerLetters,
			},
			want: errMsg,
		},
		{
			name: "验证纯特殊字符",
			args: args{
				Password: tenSpecialChar,
			},
			want: errMsg,
		},
		{
			name: "验证2种字符类型:数字+大写字母",
			args: args{
				Password: tenNumbers + tenUpperLetters,
			},
			want: errMsg,
		},
		{
			name: "验证2种字符类型:数字+小写字母",
			args: args{
				Password: tenNumbers + tenLowerLetters,
			},
			want: errMsg,
		},
		{
			name: "验证2种字符类型:数字+特殊字符",
			args: args{
				Password: tenNumbers + tenSpecialChar,
			},
			want: errMsg,
		},
		{
			name: "验证2种字符类型:大写字母+小写字母",
			args: args{
				Password: tenUpperLetters + tenLowerLetters,
			},
			want: errMsg,
		},
		{
			name: "验证2种字符类型:大写字母+特殊字符",
			args: args{
				Password: tenUpperLetters + tenSpecialChar,
			},
			want: errMsg,
		},
		{
			name: "验证2种字符类型:小写字母+特殊字符",
			args: args{
				Password: tenLowerLetters + tenSpecialChar,
			},
			want: errMsg,
		},
		{
			name: "验证3种字符类型:数字+大写字母+小写字母",
			args: args{
				Password: tenNumbers + tenUpperLetters + tenLowerLetters,
			},
			want: "",
		},
		{
			name: "验证3种字符类型:数字+大写字母+特殊字符",
			args: args{
				Password: tenNumbers + tenUpperLetters + tenSpecialChar,
			},
			want: "",
		},
		{
			name: "验证3种字符类型:数字+小写字母+特殊字符",
			args: args{
				Password: tenNumbers + tenLowerLetters + tenSpecialChar,
			},
			want: "",
		},
		{
			name: "验证3种字符类型:大写字母+小写字母+特殊字符",
			args: args{
				Password: tenUpperLetters + tenLowerLetters + tenSpecialChar,
			},
			want: "",
		},
	}
	type ValidateCase struct {
		Password string `validate:"min=8,max=30,charTypeGT3"`
	}
	v, err := testValidator()
	if err != nil {
		t.Fatalf("init validator failed.err:[%s]", err.Error())
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ValidateCase{Password: tt.args.Password}
			err := v.Struct(c)
			if got := Translate2Chinese(err); got != tt.want {
				t.Errorf("allchar failed. got:[%s] want:[%s]", got, tt.want)
			}
		})
	}
}

func Test_validateOneOfMembers(t *testing.T) {
	type args struct {
		Provider string
	}
	errMsg := fmt.Sprintf(mustInTransErr, getMustInErrMsg(mustInCloudParam))
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "一个部署于 cloud 的成员",
			args: args{
				Provider: "TencentCloud",
			},
			want: fmt.Sprintf("[Password] %s", errMsg),
		},
	}
	type ValidateCase struct {
		Provider string `validate:"mustIn=cloud"`
	}
	v, err := testValidator()
	if err != nil {
		t.Fatalf("init validator failed.err:[%s]", err.Error())
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ValidateCase{Provider: tt.args.Provider}
			err := v.Struct(c)
			if got := Translate2Chinese(err); got != tt.want {
				t.Errorf("allchar failed. got:[%s] want:[%s]", got, tt.want)
			}
		})
	}
	return
}
