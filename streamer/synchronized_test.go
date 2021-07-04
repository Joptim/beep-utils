package streamer

import (
	"github.com/joptim/go-mediaplayer/buffers"
	"math"
	"os"
	"testing"
)

var buf buffers.IBuffers

func TestMain(m *testing.M) {
	// Load buffers
	buf = buffers.New(nil)
	paths := []string{
		"../testdata/trackA.mp3",
		"../testdata/trackB.mp3",
	}
	for _, path := range paths {
		if err := buf.Load(path); err != nil {
			panic(err)
		}
	}

	// Run tests
	code := m.Run()

	// Teardown
	buf.ReleaseAll()

	os.Exit(code)
}

func TestMixer_Add(t *testing.T) {
	ss := buffers.GetStreamSeeker(buf, "../testdata/trackA.mp3", t)
	mixer := NewMixer()
	id := "foo"
	err := mixer.Add(ss, id)
	if err != nil {
		t.Fatalf("with id %s, expected nil error, got non-nil error", id)
	}

	// Assert Synchronized is initialised
	if !mixer.isInitialized {
		t.Errorf("with id %s, expected mixer to be initialised", id)
	}

	// Assert length is updated
	if mixer.Len() != ss.Len() {
		t.Errorf(
			"with id %s, expected length %d, got length %d",
			id,
			mixer.Len(),
			ss.Len(),
		)
	}

	// Assert position defaults to zero
	if 0 != mixer.Position() {
		t.Errorf(
			"with id %s, expected position %d, got position %d",
			id,
			0,
			mixer.Position(),
		)
	}
}

func TestMixer_Add_ReturnsErrorIfStreamerLengthDoesNotMatch(t *testing.T) {
	path := "../testdata/trackA.mp3"
	ss := buffers.GetStreamSeeker(buf, path, t)

	mixer := NewMixer()
	AddStreamSeeker(ss, "foo", mixer, t)

	// Simulate length of previous streamer is different (hacky)
	mixer.length = mixer.Len() - 1

	// Add a streamer with non-matching length
	id := "bar"
	if err := mixer.Add(ss, id); err == nil {
		t.Errorf("with id %s, expected non-nil error, got nil error", id)
	}
}

func TestMixer_Remove(t *testing.T) {
	ss := buffers.GetStreamSeeker(buf, "../testdata/trackA.mp3", t)
	mixer := NewMixer()
	id := "foo"
	AddStreamSeeker(ss, id, mixer, t)
	mixer.Remove(id)
	if _, exists := mixer.streamers[id]; exists {
		t.Errorf("with id %s, streamer not expected", id)
	}
}

func TestMixer_Seek(t *testing.T) {
	path := "../testdata/trackA.mp3"
	ss := buffers.GetStreamSeeker(buf, path, t)

	mixer := NewMixer()
	AddStreamSeeker(ss, "foo", mixer, t)
	position := 10
	if err := mixer.Seek(position); err != nil {
		t.Fatalf("with position %d, expected nil error, got %v", position, err)
	}

	if position != mixer.Position() {
		t.Fatalf(
			"with position %d, expected position %d, got %d",
			position,
			position,
			mixer.Position(),
		)
	}

}

func TestMixer_Seek_ReturnsErrorIfOutOfBounds(t *testing.T) {
	path := "../testdata/trackA.mp3"
	ss := buffers.GetStreamSeeker(buf, path, t)

	mixer := NewMixer()
	AddStreamSeeker(ss, "foo", mixer, t)

	positions := []int{-1, mixer.Len() + 1}
	for _, p := range positions {
		if err := mixer.Seek(p); err == nil {
			t.Fatalf("with position %d, expected non-nil error, got nil error", p)
		}
	}
}

func TestMixer_Position_ReturnsCurrentPositionAfterStream(t *testing.T) {
	path := "../testdata/trackA.mp3"
	ss := buffers.GetStreamSeeker(buf, path, t)

	mixer := NewMixer()
	AddStreamSeeker(ss, "foo", mixer, t)

	length := 10
	sample := make([][2]float64, length, length)
	mixer.Stream(sample)
	if length != mixer.Position() {
		t.Fatalf(
			"with length %d, expected position %d, got %d",
			length,
			length,
			mixer.Position(),
		)
	}
}

func TestMixer_Stream(t *testing.T) {
	path := "../testdata/trackA.mp3"
	ss := buffers.GetStreamSeeker(buf, path, t)

	mixer := NewMixer()
	AddStreamSeeker(ss, "foo", mixer, t)

	Seek(3000, mixer, t)
	length := 10
	sample := make([][2]float64, length, length)
	n, ok := mixer.Stream(sample)

	if length != n {
		t.Fatalf("expected Stream to return %d samples, got %d", length, n)
	}

	if !ok {
		t.Fatalf("expected Stream to return %t, got %t", true, ok)
	}

	// Assert sample is filled with values
	for _, sp := range sample {
		if math.Abs(sp[0]) <= 1e-12 && math.Abs(sp[1]) <= 1e-12 {
			t.Fatalf("expected stream to be filled with non-zero values, got non-zero values")
		}
	}
}

func TestMixer_Stream_FillsShorterSample(t *testing.T) {
	path := "../testdata/trackA.mp3"
	ss := buffers.GetStreamSeeker(buf, path, t)

	mixer := NewMixer()
	AddStreamSeeker(ss, "foo", mixer, t)

	size := 5
	Seek(ss.Len()-size, mixer, t)
	sample := make([][2]float64, 2*size, 2*size)
	n, ok := mixer.Stream(sample)

	if size != n {
		t.Fatalf("expected Stream to return %d samples, got %d", size, n)
	}

	if ok {
		t.Fatalf("expected Stream to return %t, got %t", false, ok)
	}
}

func TestMixer_Stream_Drained(t *testing.T) {
	path := "../testdata/trackA.mp3"
	ss := buffers.GetStreamSeeker(buf, path, t)

	mixer := NewMixer()
	AddStreamSeeker(ss, "foo", mixer, t)

	Seek(mixer.Len(), mixer, t)
	sample := make([][2]float64, 10, 10)
	n, ok := mixer.Stream(sample)
	expectedSamples := 0

	if expectedSamples != n {
		t.Fatalf("expected Stream to return %d samples, got %d", expectedSamples, n)
	}

	if ok {
		t.Fatalf("expected Stream to return %t, got %t", false, ok)
	}
}
