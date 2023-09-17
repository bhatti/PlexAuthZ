package domain

// Factory helper
type Factory[T any] func() *T

// Closeable can be closed
type Closeable interface {
	Close() error
}
