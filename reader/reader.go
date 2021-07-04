package reader

import (
	"bytes"
	"io"
	"io/ioutil"
)

type FileReader struct{}

func (f FileReader) GetReadCloser(path string) (io.ReadCloser, error) {
	if bytes_, err := ioutil.ReadFile(path); err != nil {
		return nil, err
	} else {
		return io.NopCloser(bytes.NewReader(bytes_)), nil
	}
}

type DummyReader struct {
	Contents []byte
}

func (d DummyReader) GetReadCloser(_ string) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(d.Contents)), nil
}
