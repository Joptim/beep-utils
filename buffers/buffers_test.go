package buffers

import (
	"github.com/joptim/beep-utils/reader"
	"testing"
)

func TestNew_UsesFileReaderByDefault(t *testing.T) {
	buf, ok := New(nil).(*Buffers)
	if !ok {
		t.Error("failed casting interface to concrete type")
	}
	_, isFileReader := buf.reader.(reader.FileReader)
	if !isFileReader {
		t.Errorf("with nil argument, expected reader.FileReader, got %T", buf.reader)
	}
}

func TestBuffers_Load(t *testing.T) {
	buf := New(nil)
	path := "../testdata/trackA.mp3"
	err := buf.Load(path)
	if err != nil {
		t.Fatalf("with %s, got error %v, expected nil error", path, err)
	}
	buffType := CastToBuffers(buf, t)

	// Assert buffer is loaded
	if _, exists := buffType.buffers[path]; !exists {
		t.Fatalf("with %s, buffer not loaded, expected buffer to be loaded", path)
	}

	// Assert is initialised
	if !buffType.isInitialised {
		t.Fatalf("with %s, object not initialised, expected object to be initialised", path)
	}
}

func TestBuffers_Load_ReturnsErrorIfFileDoesNotExist(t *testing.T) {
	buf := New(nil)
	path := "/foo/bar/baz.mp3"
	if err := buf.Load(path); err == nil {
		t.Errorf("with path %v, got nil error, expected non-nil error", path)
	}
}

func TestBuffers_Load_ReturnsErrorIfFileIsNotMp3(t *testing.T) {
	dummy := reader.DummyReader{Contents: []byte("foo bar baz qux")}
	buf := New(dummy)
	if err := buf.Load("/dummy/path.mp3"); err == nil {
		t.Errorf("with %T, got nil error, expected non-nil error", dummy)
	}
}

func TestBuffers_Release(t *testing.T) {
	buf := New(nil)
	path := "../testdata/trackA.mp3"
	LoadPath(buf, path, t)
	buf.Release(path)

	buffType := CastToBuffers(buf, t)
	// Assert buffer is not loaded
	if _, exists := buffType.buffers[path]; exists {
		t.Fatalf("with %s, buffer not expected", path)
	}
}

func TestBuffers_ReleaseAll(t *testing.T) {
	buf := New(nil)
	pathA := "../testdata/trackA.mp3"
	pathB := "../testdata/trackB.mp3"
	LoadPath(buf, pathA, t)
	LoadPath(buf, pathB, t)
	buf.ReleaseAll()

	buffType := CastToBuffers(buf, t)
	// Assert buffer is not loaded
	if _, exists := buffType.buffers[pathA]; exists {
		t.Fatalf("with %s, buffer not expected", pathA)
	}
	if _, exists := buffType.buffers[pathB]; exists {
		t.Fatalf("with %s, buffer not expected", pathB)
	}
}

func TestBuffers_GetFormat_ReturnsErrorIfNotInitialised(t *testing.T) {
	buf := New(nil)
	if _, err := buf.GetFormat(); err == nil {
		t.Errorf("expected non-nil error, got nil error")
	}
}
