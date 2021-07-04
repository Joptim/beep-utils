package buffers

import (
	"github.com/faiface/beep"
	"testing"
)

func CastToBuffers(b IBuffers, t *testing.T) *Buffers {
	t.Helper()
	casted, ok := b.(*Buffers)
	if !ok {
		t.Fatalf("cannot cast %T to Buffers type", b)
	}
	return casted
}

func LoadPath(b IBuffers, path string, t *testing.T) {
	t.Helper()
	if err := b.Load(path); err != nil {
		t.Fatalf("with %s, got error %v, expected nil error", path, err)
	}
}

func GetStreamSeeker(b IBuffers, path string, t *testing.T) beep.StreamSeeker {
	t.Helper()
	streamSeeker, err := b.GetStreamSeeker(path)
	if err != nil {
		t.Fatalf("cannot GetStreamSeeker from %s: %v", path, err)
	}
	return streamSeeker
}
