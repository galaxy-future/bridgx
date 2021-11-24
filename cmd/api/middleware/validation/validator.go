package validation

import (
	"errors"

	chinese "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
)

type Validation struct {
	validateFunc     func(validator.FieldLevel) bool
	translateFunc    func(ut ut.Translator, fe validator.FieldError) string
	translateRegFunc func(ut ut.Translator) error
	tag              string
}

type TagValidatorMap map[string]Validation

// appendTagValidation not thread safe and validation with same tag name will be replaced.
func appendTagValidation(tag string, validation Validation) {
	if len(tagMap) == 0 {
		tagMap = make(map[string]Validation)
	}
	tagMap[tag] = validation
}

// appendMultiTagValidation not thread safe.
func appendMultiTagValidation(validations ...Validation) {
	if len(validations) == 0 {
		return
	}
	for _, v := range validations {
		appendTagValidation(v.tag, v)
	}
}

var (
	zhTranslator            = initZhTranslator()
	defaultTranslateFunc    = func(ut ut.Translator, fe validator.FieldError) string { return "参数有误" }
	defaultTranslateRegFunc = func(ut ut.Translator) error { return nil }

	tagMap = make(map[string]Validation)
)

func initZhTranslator() ut.Translator {
	uni := ut.New(chinese.New())
	trans, _ := uni.GetTranslator("zh")
	return trans
}

func RegisterValidators(v *validator.Validate) error {
	err := registerZHTranslator(v)
	if err != nil {
		return err
	}
	err = registerValidationAndTranslation(v)
	if err != nil {
		return err
	}
	return nil
}

func registerValidationAndTranslation(v *validator.Validate) error {
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
