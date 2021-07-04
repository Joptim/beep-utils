package reader

import "io"

type IReader interface {
	GetReadCloser(path string) (io.ReadCloser, error)
}
