package utils

import (
	"github.com/go-playground/validator/v10"
)

func ValidateReqPayload(reqPayload interface{}) error {
	myValidator := validator.New()
	return myValidator.Struct(reqPayload)
	//if  err != nil {
	//	return translateErrTag(err.)
	//}
	//return nil
}

//func translateErrTag(err validator.FieldError) error {
//	switch err.Tag() {
//	case "required":
//		return fmt.Errorf("The field '%s' is required", err.Field())
//	case "gte":
//		return fmt.Errorf("The field '%s' must be greater than or equal to %s", err.Field(), err.Param())
//	default:
//		return fmt.Errorf("The field '%s' is not valid", err.Field())
//	}
//}