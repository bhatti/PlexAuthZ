package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldConvertArrayToMap(t *testing.T) {
	require.Equal(t, 1, len(ArrayToMap("k", "v")))
}

func Test_ShouldConvertToInt64(t *testing.T) {
	require.Equal(t, int64(0), ToInt64(nil))
	require.Equal(t, int64(10), ToInt64(int64(10)))
	require.Equal(t, int64(10), ToInt64(int32(10)))
	require.Equal(t, int64(10), ToInt64(10))
	require.Equal(t, int64(10), ToInt64(uint(10)))
	require.Equal(t, int64(10), ToInt64("10"))
}

func Test_ShouldConvertToBoolean(t *testing.T) {
	require.False(t, ToBoolean(nil))
	require.True(t, ToBoolean(int64(10)))
	require.True(t, ToBoolean(int32(10)))
	require.True(t, ToBoolean(10))
	require.True(t, ToBoolean(uint(10)))
	require.True(t, ToBoolean("true"))
	require.False(t, ToBoolean("N"))
	require.False(t, ToBoolean(false))
	require.True(t, ToBoolean(true))
}

func Test_ShouldConvertToInt(t *testing.T) {
	require.Equal(t, 0, ToInt(nil))
	require.Equal(t, 10, ToInt(int64(10)))
	require.Equal(t, 10, ToInt(int32(10)))
	require.Equal(t, 10, ToInt(10))
	require.Equal(t, 10, ToInt(uint(10)))
}

func Test_ShouldConvertToFloat64(t *testing.T) {
	var f32 float32 = 10
	var f64 float64 = 10
	var i64 int64 = 10
	var u64 uint64 = 10
	require.Equal(t, float64(0), ToFloat64(nil))
	require.Equal(t, float64(10), ToFloat64(float64(10)))
	require.Equal(t, float64(10), ToFloat64(float32(10)))
	require.Equal(t, float64(10), ToFloat64(10))
	require.Equal(t, float64(10), ToFloat64(uint(10)))
	require.Equal(t, float64(10), ToFloat64(int32(10)))
	require.Equal(t, float64(10), ToFloat64(int64(10)))
	require.Equal(t, float64(10), ToFloat64(uint64(10)))
	require.Equal(t, float64(10), ToFloat64(&f32))
	require.Equal(t, float64(10), ToFloat64(&f64))
	require.Equal(t, float64(10), ToFloat64(&i64))
	require.Equal(t, float64(10), ToFloat64(&u64))
}
