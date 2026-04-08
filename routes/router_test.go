package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"pismo-assignment/db"
	"pismo-assignment/routes"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "test.sqlite")
	require.NoError(t, db.Connect(dsn))
	require.NoError(t, db.Migrate())
	gin.SetMode(gin.TestMode)
	return routes.SetupRouter()
}

// Accounts: create and read balance for a new account.
func TestAccounts_CreateAndGet(t *testing.T) {
	r := setupTestRouter(t)

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

// Installments: type-2 creates a plan and debits first EMI; POST .../next debits the next EMI and updates balance.
func TestInstallments_FirstEMIAndNext(t *testing.T) {
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

	inst := map[string]any{"account_id": accountID, "operation_type_id": 2, "amount": 300.0, "tenure": 3}
	ib, _ := json.Marshal(inst)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(ib))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusCreated, w2.Code)

	var txResp map[string]any
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &txResp))
	require.Equal(t, -100.0, txResp["amount"])
	plan := txResp["installment_plan"].(map[string]any)
	planID := uint(plan["plan_id"].(float64))
	require.Equal(t, float64(1), plan["paid_emis"])
	require.Equal(t, float64(2), plan["remaining_emis"])

	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodPost,
		"/accounts/"+strconv.FormatUint(uint64(accountID), 10)+"/installments/"+strconv.FormatUint(uint64(planID), 10)+"/next", nil)
	r.ServeHTTP(w3, req3)
	require.Equal(t, http.StatusOK, w3.Code)

	var payResp map[string]any
	require.NoError(t, json.Unmarshal(w3.Body.Bytes(), &payResp))
	require.Equal(t, float64(2), payResp["paid_emis"])
	require.Equal(t, float64(1), payResp["remaining_emis"])

	wBal := httptest.NewRecorder()
	reqBal := httptest.NewRequest(http.MethodGet, "/accounts/"+strconv.FormatUint(uint64(accountID), 10), nil)
	r.ServeHTTP(wBal, reqBal)
	var bal map[string]any
	require.NoError(t, json.Unmarshal(wBal.Body.Bytes(), &bal))
	// 300 credit - 100 first EMI - 100 second EMI
	require.Equal(t, float64(100), bal["balance"])
}
