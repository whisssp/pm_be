package utils

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

const (
	phonePatter string = "^+{0,1}0{0,1}62[0-9]+$"
)

var myValidator *validator.Validate = nil

func InitValidatorHelper() {
	myValidator = validator.New()
	registerValidation()
}

func ValidateReqPayload(reqPayload interface{}) error {
	return myValidator.Struct(reqPayload)
	//if  err != nil {
	//	return translateErrTag(err.)
	//}
	//return nil
}

func registerValidation() {
	regexValidate()
}

func regexValidate() {
	err := myValidator.RegisterValidation("phone", phoneValidate)
	if err != nil {
		return
	}
}

func phoneValidate(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	valid, err := regexp.Match(phonePatter, []byte(value))
	if err != nil {
		return false
	}
	return valid
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