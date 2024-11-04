package isunippets

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/profile"
	"net/http"
	"net/http/pprof"
)

func AddPprof(e *echo.Echo) {
	pprofGroup := e.Group("/debug/pprof")
	pprofGroup.Any("/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	pprofGroup.Any("/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	pprofGroup.Any("/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	pprofGroup.Any("/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	pprofGroup.Any("/*", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
}

type PprofOption struct {
	ProfilePath string
}

func RunPprof(options *PprofOption) interface {
	Stop()
} {
	if options == nil {
		options = &PprofOption{
			ProfilePath: ".",
		}
	}
	return profile.Start(profile.ProfilePath(options.ProfilePath))
}
