package validator

import "personal_bot/backend/app/validator"

func ValidateStruct(structToValidate interface{}) error {
	validate := validator.GetValidator()
	if err := validate.Struct(structToValidate); err != nil {
		return err
	}
	return nil
}
