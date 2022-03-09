package client_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/graillus/timelapse/internal/api"
	"github.com/graillus/timelapse/internal/api/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostFrame(t *testing.T) {
	t.Parallel()

	router := mux.NewRouter()
	api.ConfigureRoutes(router)
	api.StoragePath = "/tmp/test"

	server := httptest.NewServer(router)
	server.Client()

	c := client.New(server.URL)
	_, err := c.PostFrame("../../../test/fixtures/frames/1200.jpg")
	require.NoError(t, err)

	assert.FileExists(t, "/tmp/test/uploads/1200.jpg")

	_ = os.Remove("/tmp/test/uploads/1200.jpg")
}
