package buffers

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/joptim/go-mediaplayer/reader"
	"sync"
)

type Buffers struct {
	buffers       map[string]*beep.Buffer
	format        beep.Format
	rwmutex       sync.RWMutex
	reader        reader.IReader
	isInitialised bool
}

// Load loads audio files from paths into a buffer.
// Calling this method with files already loaded takes no action.
// This method only supports loading mp3 files, though it's
// not hard to add more audio formats.
func (h *Buffers) Load(path string) error {
	h.rwmutex.Lock()
	defer h.rwmutex.Unlock()
	if _, exists := h.buffers[path]; exists {
		// Path already loaded
		return nil
	}
	file, err := h.reader.GetReadCloser(path)
	if err != nil {
		return err
	}
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		return err
	}

	if !h.isInitialised {
		h.format = format
		h.isInitialised = true
	}

	buffer := beep.NewBuffer(h.format)
	if h.format.SampleRate == format.SampleRate {
		buffer.Append(streamer)
	} else {
		resampled := beep.Resample(4, format.SampleRate, h.format.SampleRate, streamer)
		buffer.Append(resampled)
	}
	h.buffers[path] = buffer
	if err := streamer.Close(); err != nil {
		return err
	}
	return nil
}

// Release frees audio files' memory from buffer.
// Calling this method with a track already released takes no action.
func (h *Buffers) Release(path string) {
	h.rwmutex.Lock()
	defer h.rwmutex.Unlock()
	delete(h.buffers, path)
}

// ReleaseAll frees all audio files' memory from buffer.
func (h *Buffers) ReleaseAll() {
	h.rwmutex.Lock()
	defer h.rwmutex.Unlock()
	for path := range h.buffers {
		delete(h.buffers, path)
	}
}

// GetStreamSeeker returns a StreamSeeker from a buffer.
func (h *Buffers) GetStreamSeeker(path string) (beep.StreamSeeker, error) {
	h.rwmutex.RLock()
	defer h.rwmutex.RUnlock()
	buffer, exists := h.buffers[path]
	if !exists {
		return nil, fmt.Errorf("cannot get Stream Seeker from unloaded buffer %s", path)
	}
	return buffer.Streamer(0, buffer.Len()), nil
}

// GetFormat returns the common format of all buffers
func (h *Buffers) GetFormat() (beep.Format, error) {
	h.rwmutex.RLock()
	defer h.rwmutex.RUnlock()
	if !h.isInitialised {
		return beep.Format{}, fmt.Errorf("cannot GetFormat, buffer not initialised")
	}
	return h.format, nil
}

// New returns an implementation of IBuffers. If reader is nil,
// then the default reader.FileReader is used.
func New(r reader.IReader) IBuffers {
	if r == nil {
		r = reader.FileReader{}
	}
	return &Buffers{
		buffers:       make(map[string]*beep.Buffer),
		format:        beep.Format{},
		rwmutex:       sync.RWMutex{},
		reader:        r,
		isInitialised: false,
	}
}
