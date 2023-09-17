package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldBuildValidationError(t *testing.T) {
	// GIVEN a mismatch error
	err := NewValidationError("test error")
	// THEN it should match message
	require.Error(t, err)
	require.Equal(t, "test error [EC100400]", err.Error())
	require.Equal(t, 400, ErrorToHTTPStatus(err))
}

func Test_ShouldBuildNotFoundError(t *testing.T) {
	// GIVEN a mismatch error
	err := NewNotFoundError("test error")
	// THEN it should match message
	require.Error(t, err)
	require.Equal(t, "test error [EC100404]", err.Error())
	require.Equal(t, 404, ErrorToHTTPStatus(err))
}

func Test_ShouldBuildInternalError(t *testing.T) {
	// GIVEN a mismatch error
	err := NewInternalError("test error", "code")
	// THEN it should match message
	require.Error(t, err)
	require.Equal(t, "test error [code]", err.Error())
	require.Equal(t, 500, ErrorToHTTPStatus(err))
}

func Test_ShouldBuildAuthError(t *testing.T) {
	// GIVEN a mismatch error
	err := NewAuthError("test error")
	// THEN it should match message
	require.Error(t, err)
	require.Equal(t, "test error [EC100401]", err.Error())
	require.Equal(t, 401, ErrorToHTTPStatus(err))
}

func Test_ShouldBuildDuplicateError(t *testing.T) {
	// GIVEN a mismatch error
	err := NewDuplicateError("test error")
	// THEN it should match message
	require.Error(t, err)
	require.Equal(t, "test error [EC100409]", err.Error())
	require.Equal(t, 409, ErrorToHTTPStatus(err))
}

func Test_ShouldBuildDatabaseError(t *testing.T) {
	// GIVEN a mismatch error
	err := NewDatabaseError("test error")
	// THEN it should match message
	require.Error(t, err)
	require.Equal(t, "test error [EC100510]", err.Error())
	require.Equal(t, 500, ErrorToHTTPStatus(err))
}
