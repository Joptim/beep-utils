package streamer

import (
	"github.com/faiface/beep"
	"testing"
)

func AddStreamSeeker(ss beep.StreamSeeker, id string, mixer *Synchronized, t *testing.T) {
	if err := mixer.Add(ss, id); err != nil {
		t.Fatalf("cannot load StreamSeeker with id %s into mixer: %v", id, err)
	}
}

func Seek(p int, mixer *Synchronized, t *testing.T) {
	if err := mixer.Seek(p); err != nil {
		t.Fatalf("cannot seek mixer to position %d: %v", p, err)
	}
}
