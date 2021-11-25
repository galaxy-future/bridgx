package validation

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const (
	CharTypeGT3         = "charTypeGT3"
	CharTypeGT3TransErr = "必须同时包含三项（大写字母、小写字母、数字、 ()`~!@#$%^&*_-+=|{}[]:;'<>,.?/ 中的特殊符号）"

	mustIn           = "mustIn"
	mustInTransErr   = "必须是 [%s] 之一"
	mustInCloudParam = "cloud"
	delimiter        = "、"
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

	mustInMembers = map[string]map[string]struct{}{
		mustInCloudParam: {"AlibabaCloud": {}, "HuaweiCloud": {}},
	}
	mustInErrMsgCache       = map[string]string{}
	mustInErrMsgCacheRWLock = sync.RWMutex{}
)

func RegisterCustomValidators() {
	appendMultiTagValidation(
		// Add your custom Validation here.
		Validation{
			validateFunc:     validateCharacterTypeGT3,
			translateFunc:    translateCharacterErr,
			translateRegFunc: defaultTranslateRegFunc,
			tag:              CharTypeGT3,
		},
		Validation{
			validateFunc:     validateMustIn,
			translateFunc:    translateMustIn,
			translateRegFunc: defaultTranslateRegFunc,
			tag:              mustIn,
		},
	)
}

func getStructFieldName(fe validator.FieldError) string {
	f := strings.Split(fe.StructNamespace(), ".")
	return f[len(f)-1]
}

// wrapErrWithStructFieldName will wrap msg with "[`StructFieldName`] "
func wrapErrWithStructFieldName(fe validator.FieldError, msg string) string {
	return fmt.Sprintf("[%s] %s", getStructFieldName(fe), msg)
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
	return wrapErrWithStructFieldName(fe, CharTypeGT3TransErr)
}

func validateMustIn(fl validator.FieldLevel) bool {
	param := fl.Param()
	members, ok := mustInMembers[param]
	if !ok {
		return true
	}
	if len(members) > 0 {
		_, ok = members[fl.Field().String()]
		return ok
	}
	return true
}

func translateMustIn(ut ut.Translator, fe validator.FieldError) string {
	return wrapErrWithStructFieldName(fe, fmt.Sprintf(mustInTransErr, getMustInErrMsg(fe.Param())))
}

func getMustInErrMsg(param string) string {
	mustInErrMsgCacheRWLock.RLock()
	msg, ok := mustInErrMsgCache[param]
	mustInErrMsgCacheRWLock.RUnlock()
	if !ok {
		mustInErrMsgCacheRWLock.Lock()
		members := make([]string, 0)
		for mem := range mustInMembers[param] {
			members = append(members, mem)
		}
		sort.Slice(members, func(i, j int) bool {
			return members[i] < members[j]
		})
		msg = strings.Join(members, delimiter)
		mustInErrMsgCache[param] = msg
		mustInErrMsgCacheRWLock.Unlock()
	}
	return msg
}
