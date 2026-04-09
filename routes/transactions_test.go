package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// Transactions: after a credit, a purchase returns the debit as a negative amount in JSON (sign convention).
func TestTransactions_CreditThenDebitResponseAmount(t *testing.T) {
	r := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`{"document_number":"99999999999"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var created map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
	accountID := uint(created["account_id"].(float64))

	credit := map[string]any{"account_id": accountID, "operation_type_id": 4, "amount": 200.0}
	cb, _ := json.Marshal(credit)
	wc := httptest.NewRecorder()
	reqc := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(cb))
	reqc.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wc, reqc)
	require.Equal(t, http.StatusCreated, wc.Code)

	purchase := map[string]any{"account_id": accountID, "operation_type_id": 1, "amount": 123.45}
	pb, _ := json.Marshal(purchase)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(pb))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusCreated, w2.Code)

	var txResp map[string]any
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &txResp))
	require.Equal(t, -123.45, txResp["amount"])
}

// Transactions: cannot debit more than available balance.
func TestTransactions_InsufficientBalance(t *testing.T) {
	r := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`{"document_number":"88888888888"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
	var created map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
	accountID := uint(created["account_id"].(float64))

	purchase := map[string]any{"account_id": accountID, "operation_type_id": 1, "amount": 10.0}
	pb, _ := json.Marshal(purchase)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(pb))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusBadRequest, w2.Code)
	var errBody map[string]any
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &errBody))
	require.Equal(t, "insufficient balance", errBody["error"])
}

// Type 2 (installment purchase): full amount debited in one transaction, same sign convention as type 1.
func TestTransactions_Type2FullDebit(t *testing.T) {
	r := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`{"document_number":"22233344455"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
	var created map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
	accountID := uint(created["account_id"].(float64))

	credit := map[string]any{"account_id": accountID, "operation_type_id": 4, "amount": 300.0}
	cb, _ := json.Marshal(credit)
	wc := httptest.NewRecorder()
	reqc := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(cb))
	reqc.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wc, reqc)
	require.Equal(t, http.StatusCreated, wc.Code)

	inst := map[string]any{"account_id": accountID, "operation_type_id": 2, "amount": 150.0}
	ib, _ := json.Marshal(inst)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(ib))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusCreated, w2.Code)

	var txResp map[string]any
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &txResp))
	require.Equal(t, float64(2), txResp["operation_type_id"])
	require.Equal(t, -150.0, txResp["amount"])

	wBal := httptest.NewRecorder()
	reqBal := httptest.NewRequest(http.MethodGet, "/accounts/"+strconv.FormatUint(uint64(accountID), 10), nil)
	r.ServeHTTP(wBal, reqBal)
	var bal map[string]any
	require.NoError(t, json.Unmarshal(wBal.Body.Bytes(), &bal))
	require.Equal(t, float64(150), bal["balance"])
}
