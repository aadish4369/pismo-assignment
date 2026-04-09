package services

import (
	"testing"

	"github.com/stretchr/testify/require"

	"pismo-assignment/models"
)

func TestTransactionService_Create_InvalidOperationTypeID(t *testing.T) {
	svc := NewTransactionService(nil, nil)
	_, err := svc.Create(1, models.OperationType(0), 100)
	require.Error(t, err)
	require.Equal(t, "invalid operation_type_id", err.Error())
}
