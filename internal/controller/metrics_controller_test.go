package controller

import (
	"github.com/bhatti/PlexAuthZ/internal/web"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ShouldCreateMetricsController(t *testing.T) {
	webServer := web.NewStubWebServer()
	_, err := NewMetricsController(webServer)
	require.NoError(t, err)
}
