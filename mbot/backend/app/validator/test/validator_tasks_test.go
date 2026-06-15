package validator_test

import (
	"errors"
	"personal_bot/backend/app/validator"
	"personal_bot/backend/internal/core/tasks"
	"personal_bot/backend/internal/core/wallets"
	"personal_bot/backend/internal/solana/utils"
	"reflect"
	"strings"
	"testing"

	"slices"

	"github.com/gagliardetto/solana-go"
	v "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func createTestWallet(t *testing.T) (*solana.PrivateKey, solana.PublicKey) {
	t.Helper()
	privKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate random private key: %v", err)
	}
	return &privKey, privKey.PublicKey()
}

func createSolanaWallet(t *testing.T) wallets.SolanaWallet {
	t.Helper()
	privKey, publicKey := createTestWallet(t)
	return wallets.SolanaWallet{
		Id:         uuid.NewString(),
		WalletName: "Main",
		PrivateKey: *privKey,
		PublicKey:  publicKey,
	}

}

func createValidBuyTask(wallet wallets.SolanaWallet, tokenAddress solana.PublicKey) *tasks.BuyTask {

	return tasks.NewBuyTask(wallet, tokenAddress,
		[]tasks.Option{
			tasks.WithComputeUnits(10000),
			tasks.WithSlippage(0.4),
		},
		[]tasks.BuyOption{
			tasks.WithBuyAmount(utils.ConvertSolToLamport(0.0001)),
			tasks.WithBuyFee(0.001),
		},
	)
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

	var validationErrors v.ValidationErrors
	ok := errors.As(errs, &validationErrors)

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
	wallet := createSolanaWallet(t)
	buyTask := createValidBuyTask(wallet, wallet.PublicKey)

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
	wallet := createSolanaWallet(t)

	buyTask := tasks.NewBuyTask(wallet, wallet.PublicKey,
		[]tasks.Option{
			tasks.WithComputeUnits(1000),
			tasks.WithSlippage(-1),
		},
		[]tasks.BuyOption{
			tasks.WithBuyAmount(utils.ConvertSolToLamport(-0.001)),
			tasks.WithBuyFee(-1),
		})

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
	wallet := createSolanaWallet(t)
	buyTask := tasks.NewBuyTask(wallet, wallet.PublicKey,
		[]tasks.Option{
			tasks.WithComputeUnits(0),
			tasks.WithSlippage(0),
		},
		[]tasks.BuyOption{
			tasks.WithBuyAmount(utils.ConvertSolToLamport(0)),
			tasks.WithBuyFee(0),
		})

	errs := validate.Struct(buyTask)
	// Collect errors for both required + gtZero
	requiredSlice := getValidationErrorSlice(t, errs, "required")
	gtZeroSlice := getValidationErrorSlice(t, errs, "gtZero")
	validationSlice := append(requiredSlice, gtZeroSlice...)

	fields := []string{"BuyAmount", "Slippage", "Fee", "ComputeUnits"}

	for _, field := range fields {
		if !slices.Contains(validationSlice, field) {
			t.Errorf("Unexpected validation error for field %s", field)
		}
	}
}

func TestFieldsWithGreaterThanOneReturnsError(t *testing.T) {
	validate := validator.GetValidator()
	wallet := createSolanaWallet(t)
	buyTask := tasks.NewBuyTask(wallet, wallet.PublicKey,
		[]tasks.Option{
			tasks.WithComputeUnits(1000),
			tasks.WithSlippage(1.2),
		},
		[]tasks.BuyOption{
			tasks.WithBuyAmount(utils.ConvertSolToLamport(1)),
			tasks.WithBuyFee(0.5),
		})

	errs := validate.Struct(buyTask)
	validationSlice := getValidationErrorSlice(t, errs, "lt")
	fields := getAllFieldsFor(buyTask, "lt")

	for _, field := range fields {
		if !slices.Contains(validationSlice, field) {
			t.Errorf("Unexpected validation error for field: %s", field)
		}
	}
}
