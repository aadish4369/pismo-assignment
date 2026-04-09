package services

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"pismo-assignment/db"
	"pismo-assignment/repository"
)

func TestAccountService_Create_EmptyDocumentNumber(t *testing.T) {
	svc := NewAccountService(nil)
	_, err := svc.Create("")
	require.Error(t, err)
	require.Equal(t, "document_number is required", err.Error())
}

func TestAccountService_GetByID_NotFound(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "test.sqlite")
	require.NoError(t, db.Connect(dsn))
	require.NoError(t, db.Migrate())

	svc := NewAccountService(repository.NewAccountRepository())
	_, err := svc.GetByID(99999)
	require.Error(t, err)
}

