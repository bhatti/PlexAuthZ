package domain

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ShouldValidateConfig(t *testing.T) {
	cfg, err := NewConfig("")
	require.NoError(t, err)
	require.NotNil(t, cfg)
}
