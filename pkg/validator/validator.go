package validator

import (
	pkgerror "backend_test/pkg/error"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"reflect"
	"strings"
)

import (
	"errors"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/shopspring/decimal"
)

type RequestValidator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func (v *RequestValidator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err != nil {
		arr := []string{}
		validatorErrs := err.(validator.ValidationErrors)
		for _, e := range validatorErrs {
			arr = append(arr, e.Translate(v.translator))
		}
		return errors.New(strings.Join(arr, ", "))
	}
	return nil
}

func New(v *validator.Validate) *RequestValidator {
	t := registerTranslation(v)
	return &RequestValidator{
		validator:  v,
		translator: t,
	}
}

func registerTranslation(v *validator.Validate) ut.Translator {
	english := en.New()
	uni := ut.New(english, english)
	t, found := uni.GetTranslator("en")
	if !found {
		log.Warn("Validation translation not found, will use default message format")
	}
	err := en_translations.RegisterDefaultTranslations(v, t)
	if err != nil {
		log.Warn("Register default translation error:", err)
	}
	registerTranslationNotBlank(v, t)
	return t
}

func registerTranslationNotBlank(v *validator.Validate, t ut.Translator) {
	err := v.RegisterTranslation("notblank", t,
		func(ut ut.Translator) error {
			return ut.Add("notblank", "{0} must not be empty or contains only whitespace characters", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("notblank", fe.Field())
			return t
		},
	)
	if err != nil {
		log.Warn("Register 'notblank' translation error:", err)
	}
}

func BindAndValidate(ctx echo.Context, req interface{}) pkgerror.CustomError {
	if err := ctx.Bind(req); err != nil {
		return pkgerror.ErrInvalidParams.WithError(err)
	}
	err := ctx.Validate(req)
	if err != nil {
		return pkgerror.ErrInvalidParams.WithError(err)
	}
	return pkgerror.NoError
}

func DecimalValidator(field reflect.Value) interface{} {
	if dec, ok := field.Interface().(decimal.Decimal); ok {
		return dec.InexactFloat64()
	}
	return nil
}
