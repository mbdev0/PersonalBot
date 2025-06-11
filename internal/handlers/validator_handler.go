package handlers

import "pump_fun/internal/launch/validator"

func ValidateStruct(structToValidate interface{}) error {
	validate := validator.GetValidator()
	if err := validate.Struct(structToValidate); err != nil {
		return err
	}
	return nil
}
