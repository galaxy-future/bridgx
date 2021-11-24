package validation

import (
	"errors"

	chinese "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
)

const (
	CharTypeGT3         = "charTypeGT3"
	CharTypeGT3TransErr = "必须同时包含三项（大写字母、小写字母、数字、 ()`~!@#$%^&*_-+=|{}[]:;'<>,.?/ 中的特殊符号）"
)

var (
	zhTranslator            = initZhTranslator()
	defaultTranslateFunc    = func(ut ut.Translator, fe validator.FieldError) string { return "参数有误" }
	defaultTranslateRegFunc = func(ut ut.Translator) error { return nil }
)

type Validation struct {
	validateFunc     func(validator.FieldLevel) bool
	translateFunc    func(ut ut.Translator, fe validator.FieldError) string
	translateRegFunc func(ut ut.Translator) error
}

var (
	numberMap      = map[byte]struct{}{'1': {}, '2': {}, '3': {}, '4': {}, '5': {}, '6': {}, '7': {}, '8': {}, '9': {}, '0': {}}
	upperLetterMap = map[byte]struct{}{'A': {}, 'B': {}, 'C': {}, 'D': {}, 'E': {}, 'F': {}, 'G': {}, 'H': {}, 'I': {}, 'J': {},
		'K': {}, 'L': {}, 'M': {}, 'N': {}, 'O': {}, 'P': {}, 'Q': {}, 'R': {}, 'S': {}, 'T': {}, 'U': {}, 'V': {}, 'W': {}, 'X': {}, 'Y': {}, 'Z': {}}
	lowerLetterMap = map[byte]struct{}{'a': {}, 'b': {}, 'c': {}, 'd': {}, 'e': {}, 'f': {}, 'g': {}, 'h': {}, 'i': {}, 'j': {}, 'k': {}, 'l': {},
		'm': {}, 'n': {}, 'o': {}, 'p': {}, 'q': {}, 'r': {}, 's': {}, 't': {}, 'u': {}, 'v': {}, 'w': {}, 'x': {}, 'y': {}, 'z': {}}
	specialCharMap = map[byte]struct{}{
		'(': {}, ')': {}, '`': {}, '~': {}, '!': {}, '@': {}, '#': {}, '$': {}, '%': {}, '^': {}, '&': {}, '*': {}, '_': {},
		'-': {}, '+': {}, '=': {}, '|': {}, '{': {}, '}': {}, '[': {}, ']': {}, ':': {}, ';': {}, '\'': {}, '<': {}, '>': {}, ',': {}, '.': {}, '?': {}, '/': {},
	}

	tagMap = map[string]Validation{
		CharTypeGT3: {
			validateFunc:     validateCharacterTypeGT3,
			translateFunc:    translateCharacterErr,
			translateRegFunc: defaultTranslateRegFunc,
		},
	}
)

func initZhTranslator() ut.Translator {
	uni := ut.New(chinese.New())
	trans, _ := uni.GetTranslator("zh")
	return trans
}

func validateCharacterTypeGT3(fl validator.FieldLevel) bool {
	field := []byte(fl.Field().String())
	var numType, upperLetterType, loweLetterType, specialChatType int
	for _, c := range field {
		_, ok := numberMap[c]
		if ok && numType == 0 {
			numType = 1
		}
		_, ok = upperLetterMap[c]
		if ok && upperLetterType == 0 {
			upperLetterType = 1
		}
		_, ok = lowerLetterMap[c]
		if ok && loweLetterType == 0 {
			loweLetterType = 1
		}
		_, ok = specialCharMap[c]
		if ok && specialChatType == 0 {
			specialChatType = 1
		}
	}
	return numType+upperLetterType+loweLetterType+specialChatType >= 3
}

func translateCharacterErr(ut ut.Translator, fe validator.FieldError) string {
	return "必须同时包含三项（大写字母、小写字母、数字、 ()`~!@#$%^&*_-+=|{}[]:;'<>,.?/ 中的特殊符号）"
}

func RegisterTools(v *validator.Validate) error {
	err := registerZHTranslator(v)
	if err != nil {
		return err
	}
	err = registerCustomerValidation(v)
	if err != nil {
		return err
	}
	return nil
}

func registerCustomerValidation(v *validator.Validate) error {
	if v == nil {
		return errors.New("empty validator")
	}
	for tag, cv := range tagMap {
		err := v.RegisterValidation(tag, cv.validateFunc)
		if err != nil {
			return err
		}
		err = v.RegisterTranslation(tag, zhTranslator, defaultTranslateRegFunc, cv.translateFunc)
		if err != nil {
			return err
		}
	}
	return nil
}

func registerZHTranslator(v *validator.Validate) error {
	return zh.RegisterDefaultTranslations(v, zhTranslator)
}
func Translate2Chinese(err error) string {
	if err == nil {
		return ""
	}
	verr, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}
	for _, err := range verr {
		return err.Translate(zhTranslator)
	}
	return ""
}
