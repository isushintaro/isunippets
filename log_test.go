package isunippets

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowLog(t *testing.T) {
	assert := assert.New(t)

	e := echo.New()
	err := ShowLog(e, true)
	assert.NoError(err)

	assert.True(e.Debug)
	assert.Equal(log.ERROR, e.Logger.Level())

	err = ShowLog(e, false)
	assert.NoError(err)

	assert.False(e.Debug)
	assert.Equal(log.OFF, e.Logger.Level())
}
