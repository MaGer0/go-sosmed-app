package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ParseErrorMessage(err error) []string {
	var validationErrs validator.ValidationErrors
	var sintaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError

	if errors.As(err, &sintaxErr) {
		return []string{"Invalid JSON sintax"}
	} else if errors.As(err, &typeErr) {
		return []string{fmt.Sprintf("Field '%s' must be of type %s", typeErr.Field, typeErr.Type)}
	} else if errors.As(err, &validationErrs) {
		out := make([]string, len(validationErrs))

		for i, fieldErr := range validationErrs {
			out[i] = customMessage(fieldErr)
		}

		return out
	} else if err.Error() == "EOF" {
		return []string{"Unexpexted end of file"}
	}
	return []string{}
}

func customMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return fieldErr.Field() + " is required"
	default:
		return "Invalid request"
	}
}
