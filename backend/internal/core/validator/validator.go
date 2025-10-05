package validator

import "pump_fun/app/validator"

func ValidateStruct(structToValidate interface{}) error {
	validate := validator.GetValidator()
	if err := validate.Struct(structToValidate); err != nil {
		return err
	}
	return nil
}
