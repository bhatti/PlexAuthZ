package benchmark

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_Benchmark(t *testing.T) {
	fn := func() error {
		if time.Now().Nanosecond()%2 == 0 {
			time.Sleep(20 * time.Millisecond)
			return fmt.Errorf("test error")
		}
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	res := Benchmark(fn, Request{
		TPS:         10,
		Duration:    time.Second,
		Percentiles: []float64{95.0, 99.0},
	})
	t.Log(res.String())
	require.Equal(t, 1, len(res.ErrorsByType))
}
