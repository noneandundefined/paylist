package httpx

import (
	"fmt"

	"github.com/go-playground/validator"
	"paylist.server/infra/locale"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()

	_ = Validate.RegisterValidation("optional_uuid", func(fl validator.FieldLevel) bool {
		uuid := fl.Field().String()
		if uuid == "" {
			return true
		}

		return validator.New().Var(uuid, "uuid") == nil
	})
}

func ValidateMsg(tr locale.Translator, err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			switch fieldError.Tag() {
			case "required":
				return fmt.Sprintf(tr.TErr("error.field-required"), fieldError.Field())
			case "len":
				return fmt.Sprintf(tr.TErr("error.field-len"), fieldError.Field(), fieldError.Param())
			case "numeric":
				return fmt.Sprintf(tr.TErr("error.field-numeric"), fieldError.Field())
			case "uuid":
				return fmt.Sprintf(tr.TErr("error.field-uuid"), fieldError.Field())
			case "gt":
				return fmt.Sprintf(tr.TErr("error.field-gt"), fieldError.Field(), fieldError.Param())
			case "max":
				return fmt.Sprintf(tr.TErr("error.field-max"), fieldError.Field(), fieldError.Param())
			case "min":
				return fmt.Sprintf(tr.TErr("error.field-min"), fieldError.Field(), fieldError.Param())
			case "optional_uuid":
				return fmt.Sprintf(tr.TErr("error.field-optional-uuid"), fieldError.Field())
			default:
				return fmt.Sprintf(tr.TErr("error.field-validation-error"), fieldError.Field())
			}
		}
	}

	return tr.TErr("error.request-validation-error")
}
