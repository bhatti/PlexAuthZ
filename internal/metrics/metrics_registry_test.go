package metrics

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ShouldCreateMetricsRegistry(t *testing.T) {
	registry := New()
	registry.Incr("id1", "org", "1")
	registry.Duration("id2", 3, "job", "ok")
	registry.Set("id3", 3, "level", "2")
	summary := registry.Summary()
	require.True(t, len(summary) > 0)
}
