package testutil

import (
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"pismo-assignment/db"
	"pismo-assignment/routes"
)

// SetupRouter connects a fresh temp SQLite DB, runs migrations, and returns the app router (Gin test mode).
func SetupRouter(t *testing.T) *gin.Engine {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "test.sqlite")
	require.NoError(t, db.Connect(dsn))
	require.NoError(t, db.Migrate())
	gin.SetMode(gin.TestMode)
	return routes.SetupRouter()
}
