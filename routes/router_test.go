package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"pismo-assignment/db"
	"pismo-assignment/routes"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "test.sqlite")
	require.NoError(t, db.ConnectWithDSN(dsn))
	require.NoError(t, db.Migrate())
	gin.SetMode(gin.TestMode)
	return routes.SetupRouter()
}

func TestCreateAndGetAccount(t *testing.T) {
	r := setupTestRouter(t)

	body := []byte(`{"document_number":"12345678900"}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
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
}

func TestCreateTransaction_SignNormalization(t *testing.T) {
	r := setupTestRouter(t)

	// Create account first.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`{"document_number":"99999999999"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var created map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
	accountID := uint(created["account_id"].(float64))

	// Purchase/withdraw types must be stored as negative.
	purchase := map[string]any{
		"account_id":        accountID,
		"operation_type_id": 1,
		"amount":            123.45,
	}
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

func TestCreateInstallmentTransaction_AndPayEMI(t *testing.T) {
	r := setupTestRouter(t)

	// Create account.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`{"document_number":"11122233344"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
	var created map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
	accountID := uint(created["account_id"].(float64))

	// Create installment purchase transaction (operation_type_id=2).
	now := time.Now().UTC().Format("2006-01-02")
	instReq := map[string]any{
		"account_id":        accountID,
		"operation_type_id": 2,
		"amount":            300.00,
		"tenure":            3,
		"start_date":        now,
	}
	ib, _ := json.Marshal(instReq)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(ib))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusCreated, w2.Code)

	var instResp map[string]any
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &instResp))
	installmentID := uint(instResp["installment_id"].(float64))
	require.Equal(t, float64(3), instResp["remaining_emis"])

	// Pay first EMI.
	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodPost, "/installments/"+strconv.FormatUint(uint64(installmentID), 10)+"/pay", nil)
	r.ServeHTTP(w3, req3)
	require.Equal(t, http.StatusOK, w3.Code)
}
