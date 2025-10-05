package validator

import (
	"fmt"
	"math/big"
	"pump_fun/pkg/logger"
	"sync"

	"github.com/go-playground/validator/v10"
)

/*
	cmp returns
	1 if v > 0
	0 if v == 0
	-1 if v < 0
*/

var validate *validator.Validate
var once sync.Once

func GetValidator() *validator.Validate {
	once.Do(setupValidator)
	return validate
}

func setupValidator() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	err := validate.RegisterValidation("gtZero", func(fl validator.FieldLevel) bool {
		field := fl.Field()
		if !field.IsValid() {
			fmt.Printf("Field is not valid: %v\n", field)
			return false
		}

		switch v := field.Interface().(type) {
		case big.Int:
			return v.Cmp(big.NewInt(0)) > 0
		case *big.Int:
			if v == nil {
				return false
			}
			return v.Cmp(big.NewInt(0)) > 0
		default:
			return false
		}
	})

	if err != nil {
		logger.Error("error registering validation for gtZero")
		return
	}
}
