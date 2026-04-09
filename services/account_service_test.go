package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountService_Create_EmptyDocumentNumber(t *testing.T) {
	svc := NewAccountService(nil)
	_, err := svc.Create("")
	require.Error(t, err)
	require.Equal(t, "document_number is required", err.Error())
}
