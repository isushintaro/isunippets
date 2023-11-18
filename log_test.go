package isunippets

import (
	"bytes"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShowLog(t *testing.T) {
	assert := assert.New(t)

	check := func(c echo.Context) error {
		logger := c.Echo().Logger
		logger.Debug("debug log")
		logger.Info("info log")
		logger.Warn("warn log")
		logger.Error("error log")

		return c.String(http.StatusOK, "OK")
	}
	checkPath := "/check"

	setup := func() (*echo.Echo, *resty.Client, *bytes.Buffer) {
		e := echo.New()
		e.GET(checkPath, check)

		testServer := httptest.NewServer(e.Server.Handler)
		t.Cleanup(func() {
			testServer.Close()
		})

		client := resty.New().SetBaseURL(testServer.URL)

		output := bytes.NewBufferString("")
		e.Logger.SetOutput(output)

		return e, client, output
	}

	t.Run("show logs", func(t *testing.T) {
		e, client, output := setup()

		err := ShowLog(e, true)
		assert.NoError(err)

		res, err := client.R().Get(checkPath)
		assert.NoError(err)
		assert.Equal(http.StatusOK, res.StatusCode())

		assert.Contains(output.String(), "/check")
		assert.Contains(output.String(), "debug log")
		assert.Contains(output.String(), "info log")
		assert.Contains(output.String(), "warn log")
		assert.Contains(output.String(), "error log")
	})

	t.Run("hide logs", func(t *testing.T) {
		e, client, output := setup()

		err := ShowLog(e, false)
		assert.NoError(err)

		res, err := client.R().Get(checkPath)
		assert.NoError(err)
		assert.Equal(http.StatusOK, res.StatusCode())

		assert.Equal("", output.String())
	})
}
