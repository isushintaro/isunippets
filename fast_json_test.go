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

func TestFastJSON(t *testing.T) {
	type JSONResponse struct {
		Message string `json:"message"`
		Amount  int64  `json:"amount"`
	}

	assert := assert.New(t)

	check := func(c echo.Context) error {
		response := JSONResponse{
			Message: "hello",
			Amount:  100,
		}

		return FastJSON(c, http.StatusOK, response)
	}
	checkPath := "/checkFastJson"

	e := echo.New()
	e.GET(checkPath, check)

	testServer := httptest.NewServer(e.Server.Handler)
	t.Cleanup(func() {
		testServer.Close()
	})

	client := resty.New().SetBaseURL(testServer.URL)

	res, err := client.R().Get(checkPath)
	assert.NoError(err)

	assert.Equal(http.StatusOK, res.StatusCode())
	assert.JSONEq(`{"message":"hello","amount":100}`, res.String())
}

func TestFastJsonDecode(t *testing.T) {
	type JSONResponse struct {
		Message string `json:"message"`
		Amount  int64  `json:"amount"`
	}

	assert := assert.New(t)

	data := bytes.NewBufferString(`{"message":"hello","amount":100}`)
	response := JSONResponse{}

	err := FastJsonDecode(data, &response)
	assert.NoError(err)

	assert.Equal("hello", response.Message)
	assert.Equal(int64(100), response.Amount)
}

func TestFastJSONMarshal(t *testing.T) {
	type JSONResponse struct {
		Message string `json:"message"`
		Amount  int64  `json:"amount"`
	}

	assert := assert.New(t)

	response := JSONResponse{
		Message: "hello",
		Amount:  100,
	}

	data, err := FastJSONMarshal(response)
	assert.NoError(err)

	assert.JSONEq(`{"message":"hello","amount":100}`, string(data))
}

func TestFastJSONUnmarshal(t *testing.T) {
	type JSONResponse struct {
		Message string `json:"message"`
		Amount  int64  `json:"amount"`
	}

	assert := assert.New(t)

	data := []byte(`{"message":"hello","amount":100}`)
	response := JSONResponse{}

	err := FastJSONUnmarshal(data, &response)
	assert.NoError(err)

	assert.Equal("hello", response.Message)
	assert.Equal(int64(100), response.Amount)
}
