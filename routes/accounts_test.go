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
