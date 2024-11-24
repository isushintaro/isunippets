package isunippets

import (
	gojson "github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
	"io"
)

func FastJSON(c echo.Context, code int, i interface{}) error {
	b, err := FastJSONMarshal(i)
	if err != nil {
		return err
	}
	return c.JSONBlob(code, b)
}

func FastJsonDecode(r io.Reader, v any) error {
	return gojson.NewDecoder(r).Decode(v)
}

func FastJSONMarshal(i interface{}) ([]byte, error) {
	return gojson.Marshal(i)
}

func FastJSONUnmarshal(data []byte, i interface{}) error {
	return gojson.Unmarshal(data, i)
}
