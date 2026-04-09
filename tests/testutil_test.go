package tests

import (
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"pismo-assignment/db"
	"pismo-assignment/routes"
)

func setupRouter(t *testing.T) *gin.Engine {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "test.sqlite")
	require.NoError(t, db.Connect(dsn))
	require.NoError(t, db.Migrate())
	gin.SetMode(gin.TestMode)
	return routes.SetupRouter()
}
