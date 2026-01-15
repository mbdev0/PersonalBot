package validator_test

import (
	"math/big"
	"personal_bot/app/validator"
	"testing"
)

func TestGetValidator(t *testing.T) {
	validate := validator.GetValidator()
	if validate == nil {
		t.Error("Expected non-nil validator instance")
	}
}

func TestValidatorIsSingleton(t *testing.T) {
	validator1 := validator.GetValidator()
	validator2 := validator.GetValidator()

	if validator1 != validator2 {
		t.Error("Expected GetValidator to return the same instance on multiple calls")
	}
}

func TestGtZeroValidationWithBigInt(t *testing.T) {
	validate := validator.GetValidator()

	type TestValue struct {
		Val big.Int `validate:"gtZero"`
	}

	bigIntZero := TestValue{Val: *big.NewInt(0)}
	err := validate.Struct(&bigIntZero)
	if err == nil {
		t.Error("Expected validation error for big.Int == 0")
	}

	bigIntNegative := TestValue{Val: *big.NewInt(-1)}
	err = validate.Struct(&bigIntNegative)
	if err == nil {
		t.Error("Expected validation error for big.Int == 0")
	}

	bigIntOne := TestValue{Val: *big.NewInt(1)}
	err = validate.Struct(&bigIntOne)
	if err != nil {
		t.Errorf("Unexpected error for big.Int > 0: %v", err)
	}

}

func TestGtZeroWithBigIntPtr(t *testing.T) {
	type TestPtrValue struct {
		Val *big.Int `validate:"gtZero"`
	}

	validate := validator.GetValidator()
	zeroPtrStruct := TestPtrValue{Val: big.NewInt(0)}
	err := validate.Struct(&zeroPtrStruct)

	if err == nil {
		t.Error("Expected validation error for *big.Int == 0")
	}

	onePtrStruct := TestPtrValue{Val: big.NewInt(1)}
	err = validate.Struct(&onePtrStruct)
	if err != nil {
		t.Errorf("Unexpected error for *big.Int > 0: %v", err)
	}

	negativePtrStruct := TestPtrValue{Val: big.NewInt(-1)}
	err = validate.Struct(&negativePtrStruct)
	if err == nil {
		t.Error("Expected validation error for *big.Int < 0")
	}

	nilPtrStruct := TestPtrValue{Val: nil}
	err = validate.Struct(&nilPtrStruct)
	if err == nil {
		t.Error("Expected validation error for nil *big.Int")
	}
}
