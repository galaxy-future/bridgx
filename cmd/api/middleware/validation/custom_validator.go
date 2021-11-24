package validation

import (
	"fmt"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const (
	CharTypeGT3         = "charTypeGT3"
	CharTypeGT3TransErr = "必须同时包含三项（大写字母、小写字母、数字、 ()`~!@#$%^&*_-+=|{}[]:;'<>,.?/ 中的特殊符号）"
)

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
)

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
	return wrapErrWithStructFieldName(fe, CharTypeGT3TransErr)
}

// wrapErrWithStructFieldName will wrap msg with "[`StructFieldName`] "
func wrapErrWithStructFieldName(fe validator.FieldError, msg string) string {
	return fmt.Sprintf("[%s] %s", getStructFieldName(fe), msg)
}

func getStructFieldName(fe validator.FieldError) string {
	f := strings.Split(fe.StructNamespace(), ".")
	return f[len(f)-1]
}

func RegisterCustomValidators() {
	appendMultiTagValidation(
		// Add your custom Validation here.
		Validation{
			validateFunc:     validateCharacterTypeGT3,
			translateFunc:    translateCharacterErr,
			translateRegFunc: defaultTranslateRegFunc,
			tag:              CharTypeGT3,
		},
	)
}
