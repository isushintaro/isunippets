package isunippets

import (
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
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

func TestRunPprof(t *testing.T) {
	assert := assert.New(t)

	tempDir, err := os.MkdirTemp("", "TestRunPprof")
	assert.NoError(err)

	count(tempDir)

	stat, err := os.Stat(path.Join(tempDir, "cpu.pprof"))
	assert.NoError(err)
	assert.Greater(stat.Size(), int64(0))
}

func count(profilePath string) {
	option := &PprofOption{
		ProfilePath: profilePath,
	}
	defer RunPprof(option).Stop()

	summary := 0
	for i := 0; i < 1000000; i++ {
		summary += i
	}
}
