package controller

import (
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/web"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ShouldSucceedWithControllersRegistration(t *testing.T) {
	webServer := web.NewStubWebServer()
	to, _, err := newTestAuthController()
	require.NoError(t, err)
	require.NoError(t, StartControllers(to.config, to.authService, webServer))
}

func Test_ShouldSucceedWithSetupWebServerForTesting(t *testing.T) {
	_, err := domain.NewConfig("")
	require.NoError(t, err)
	//_, teardown := SetupWebServerForTesting(t, cfg, nil)
	//teardown()
}
