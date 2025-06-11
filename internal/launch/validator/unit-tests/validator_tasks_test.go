package validator_test

import (
	"pump_fun/internal/handlers"
	"pump_fun/internal/launch/validator"
	"pump_fun/internal/models/tasks"
	"reflect"
	"strings"
	"testing"

	"slices"

	"github.com/gagliardetto/solana-go"
	v "github.com/go-playground/validator/v10"
)

func createTestWallet(t *testing.T) (*solana.PrivateKey, solana.PublicKey) {
	t.Helper()
	privKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate random private key: %v", err)
	}
	return &privKey, privKey.PublicKey()
}

func createValidBuyTask(wallet *solana.PrivateKey, tokenAddress solana.PublicKey) *tasks.BuyTask {
	return &tasks.BuyTask{
		Wallet:       *wallet,
		TokenAddress: tokenAddress,
		BuyAmount:    handlers.ConvertSolToLamport(0.001),
		Slippage:     0.20,
		BuyFee:       0.0001,
		ComputeUnits: 200000,
	}
}

func getAllFieldsFor(task interface{}, tag string) []string {
	taskType := reflect.TypeOf(task)
	if taskType.Kind() == reflect.Ptr {
		taskType = taskType.Elem()
	}

	var fieldsWithTag []string
	for i := 0; i < taskType.NumField(); i++ {
		field := taskType.Field(i)
		if tagValue := field.Tag.Get("validate"); strings.Contains(tagValue, tag) {
			fieldsWithTag = append(fieldsWithTag, field.Name)
		}
	}
	return fieldsWithTag
}

func getValidationErrorSlice(t *testing.T, errs error, tag string) []string {
	t.Helper()

	validationErrors, ok := errs.(v.ValidationErrors)
	if !ok {
		t.Fatalf("Expected validation errors, got: %v", errs)
	}
	if len(validationErrors) == 0 {
		t.Error("Expected validation errors for empty BuyTask, got none")
	}

	var errorSlice []string

	for _, err := range validationErrors {
		if !strings.Contains(err.Tag(), tag) {
			continue
		}
		errorSlice = append(errorSlice, err.Field())
	}

	if len(errorSlice) == 0 {
		t.Error("Expected validation errors, got none")
	}

	return errorSlice
}

func TestValidatesBuyTaskCorrectly(t *testing.T) {
	validate := validator.GetValidator()
	wallet, tokenAddress := createTestWallet(t)
	buyTask := createValidBuyTask(wallet, tokenAddress)

	err := validate.Struct(buyTask)

	if err != nil {
		t.Errorf("Expected no validation error for valid BuyTask, got: %v", err)
	}
}

func TestNoFieldsReturnsRequiredError(t *testing.T) {
	validate := validator.GetValidator()
	buyTask := &tasks.BuyTask{}

	errs := validate.Struct(buyTask)
	validationSlice := getValidationErrorSlice(t, errs, "required")
	fields := getAllFieldsFor(buyTask, "required")

	for _, field := range fields {
		if !slices.Contains(validationSlice, field) {
			t.Errorf("Unexpected validation error for field: %s", field)
		}
	}

}

func TestFieldsWithLessThanZeroReturnsError(t *testing.T) {
	validate := validator.GetValidator()
	wallet, tokenAddress := createTestWallet(t)
	buyTask := &tasks.BuyTask{
		Wallet:       *wallet,
		TokenAddress: tokenAddress,
		BuyAmount:    handlers.ConvertSolToLamport(-0.001),
		Slippage:     -1,
		BuyFee:       -1,
		ComputeUnits: 1000,
	}

	errs := validate.Struct(buyTask)
	gtValidationSlice := getValidationErrorSlice(t, errs, "gt")
	fields := getAllFieldsFor(buyTask, "gt")

	for _, field := range fields {
		if !slices.Contains(gtValidationSlice, field) {
			t.Errorf("Unexpected validation error for field: %s", field)
		}
	}

	gtZeroValidationSlice := getValidationErrorSlice(t, errs, "gtZero")
	fields = getAllFieldsFor(buyTask, "gtZero")
	for _, field := range fields {
		if !slices.Contains(gtZeroValidationSlice, field) {
			t.Errorf("Unexpected validation error for field: %s", field)
		}
	}
}

func TestFieldsWithZeroEntriesReturnRequiredError(t *testing.T) {
	validate := validator.GetValidator()
	wallet, tokenAddress := createTestWallet(t)
	buyTask := &tasks.BuyTask{
		Wallet:       *wallet,
		TokenAddress: tokenAddress,
		BuyAmount:    handlers.ConvertSolToLamport(0),
		Slippage:     0,
		BuyFee:       0,
		ComputeUnits: 0,
	}

	errs := validate.Struct(buyTask)
	validationSlice := getValidationErrorSlice(t, errs, "required")
	fields := []string{"BuyAmount", "Slippage", "BuyFee", "ComputeUnits"}

	for _, field := range fields {
		if !slices.Contains(validationSlice, field) {
			t.Errorf("Unexpected validation error for field %s", field)
		}
	}
}

func TestFieldsWithGreaterThanOneReturnsError(t *testing.T) {
	validate := validator.GetValidator()
	wallet, tokenAddress := createTestWallet(t)
	buyTask := createValidBuyTask(wallet, tokenAddress)
	buyTask.Slippage = 1.2

	errs := validate.Struct(buyTask)
	validationSlice := getValidationErrorSlice(t, errs, "lt")
	fields := getAllFieldsFor(buyTask, "lt")

	for _, field := range fields {
		if !slices.Contains(validationSlice, field) {
			t.Errorf("Unexpected validation error for field: %s", field)
		}
	}
}
