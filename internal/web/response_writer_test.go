package web

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_StubResponseWriterShouldReadStatus(t *testing.T) {
	w := &StubResponseWriter{status: 1}
	require.Equal(t, 1, w.Status())
	w.WriteHeader(200)
	require.Equal(t, 200, w.Status())
}

func Test_StubResponseWriterShouldReadSize(t *testing.T) {
	w := &StubResponseWriter{size: 5}
	require.Equal(t, 5, w.Size())
}

func Test_StubResponseWriterShouldReadHeader(t *testing.T) {
	w := &StubResponseWriter{size: 1}
	require.Equal(t, 0, len(w.Header()))
}

func Test_StubResponseWriterShouldNotWrite(t *testing.T) {
	w := &StubResponseWriter{size: 1}
	n, err := w.Write([]byte("test"))
	require.NoError(t, err)
	require.Equal(t, 0, n)
}
