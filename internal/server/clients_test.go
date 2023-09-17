package server

import (
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ClientTypesMap(t *testing.T) {
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	m, err := ClientTypesMap(cfg)
	require.NoError(t, err)
	require.True(t, len(m) > 0)
}
