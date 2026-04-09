package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"pismo-assignment/internal/testutil"
)

func TestTransactions_CreditThenPurchase_NegativeAmountInJSON(t *testing.T) {
	r := testutil.SetupRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`{"document_number":"99999999999"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var created map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
	accountID := uint(created["account_id"].(float64))

	credit := map[string]any{"account_id": accountID, "operation_type_id": 4, "amount": 200.0}
	cb, err := json.Marshal(credit)
	require.NoError(t, err)
	wc := httptest.NewRecorder()
	reqc := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(cb))
	reqc.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wc, reqc)
	require.Equal(t, http.StatusCreated, wc.Code)

	purchase := map[string]any{"account_id": accountID, "operation_type_id": 1, "amount": 123.45}
	pb, err := json.Marshal(purchase)
	require.NoError(t, err)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(pb))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusCreated, w2.Code)

	var txResp map[string]any
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &txResp))
	require.Equal(t, -123.45, txResp["amount"])
}

func TestTransactions_InsufficientBalance(t *testing.T) {
	r := testutil.SetupRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`{"document_number":"88888888888"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var created map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
	accountID := uint(created["account_id"].(float64))

	purchase := map[string]any{"account_id": accountID, "operation_type_id": 1, "amount": 10.0}
	pb, err := json.Marshal(purchase)
	require.NoError(t, err)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(pb))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusBadRequest, w2.Code)

	var errBody map[string]any
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &errBody))
	require.Equal(t, "insufficient balance", errBody["error"])
}
