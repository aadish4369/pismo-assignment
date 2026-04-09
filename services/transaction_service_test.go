package services

import (
	"testing"

	"github.com/stretchr/testify/require"

	"pismo-assignment/models"
)

func TestNormalizedAmountInPaisa(t *testing.T) {
	cases := []struct {
		name   string
		op     models.OperationType
		input  int64
		expect int64
	}{
		{"debit_positive_becomes_negative", models.NormalPurchase, 100_00, -100_00},
		{"debit_already_negative_unchanged", models.NormalPurchase, -100_00, -100_00},
		{"credit_positive_unchanged", models.CreditVoucher, 50_00, 50_00},
		{"credit_negative_becomes_positive", models.CreditVoucher, -50_00, 50_00},
		{"withdraw_same_as_purchase", models.Withdraw, 10_00, -10_00},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizedAmountInPaisa(tc.op, tc.input)
			require.Equal(t, tc.expect, got)
		})
	}
}

func TestTransactionService_Create_InvalidOperationTypeID(t *testing.T) {
	svc := NewTransactionService(nil, nil)
	_, err := svc.Create(1, models.OperationType(0), 100)
	require.Error(t, err)
	require.Equal(t, "invalid operation_type_id", err.Error())

	_, err = svc.Create(1, models.OperationType(5), 100)
	require.Error(t, err)
	require.Equal(t, "invalid operation_type_id", err.Error())
}
