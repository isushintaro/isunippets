package isunippets

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func ShowLog(e *echo.Echo, show bool) error {
	if show {
		e.Debug = true
		e.Logger.SetLevel(log.ERROR)
		e.Use(middleware.Logger())
	} else {
		e.Debug = false
		e.Logger.SetLevel(log.OFF)
	}

	return nil
}
