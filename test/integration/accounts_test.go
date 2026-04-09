package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"pismo-assignment/internal/testutil"
)

func TestAccounts_CreateAndGetBalance(t *testing.T) {
	r := testutil.SetupRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`{"document_number":"12345678900"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var created map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
	id := uint(created["account_id"].(float64))

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/accounts/"+strconv.FormatUint(uint64(id), 10), nil)
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusOK, w2.Code)

	var got map[string]any
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &got))
	require.Equal(t, float64(0), got["balance"])
}

func TestAccounts_Get_InvalidAccountID(t *testing.T) {
	r := testutil.SetupRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/accounts/not-a-number", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
