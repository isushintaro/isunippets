package isunippets

import (
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddPprof(t *testing.T) {
	assert := assert.New(t)

	e := echo.New()

	testServer := httptest.NewServer(e.Server.Handler)
	t.Cleanup(func() {
		testServer.Close()
	})

	AddPprof(e)

	client := resty.New().SetBaseURL(testServer.URL)
	res, err := client.R().Get("/debug/pprof/cmdline")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode())
}
