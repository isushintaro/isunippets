package isunippets

import (
	gojson "github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
)

func FastJSON(c echo.Context, code int, i interface{}) error {
	b, _ := gojson.Marshal(i)
	return c.JSONBlob(code, b)
}
