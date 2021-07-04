package streamer

import (
	"fmt"
	"github.com/faiface/beep"
	"sync"
)

// Synchronized streams multiple beep.StreamSeeker synchronously, in a similar
// fashion as beep.Mixer. Synchronized provides methods to add and remove
// streamers concurrently while streaming. All streams must have the same length.
// Synchronized implements beep.StreamSeeker.
type Synchronized struct {
	streamers map[string]*beep.StreamSeeker
	// Current position of the mixer
	position int
	// length of the streamers. All streamers must have the same length.
	length        int
	err           error
	mutex         sync.RWMutex
	isInitialized bool
}

// Len returns the total number of samples of the Streamer.
func (m *Synchronized) Len() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.length
}

// Position returns the current position of the Streamer.
// This value is between 0 and the total length.
func (m *Synchronized) Position() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.position
}

// Seek sets the position of the Streamer to the provided value.
// If an error occurs during seeking, the position remains unchanged.
func (m *Synchronized) Seek(p int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if p < 0 || p > m.length {
		return fmt.Errorf(
			"cannot seek to position %d. Position out of bounds %d, %d", p, 0, m.length,
		)
	}

	m.position = p
	return nil
}

// Stream copies at most len(samples) next audio samples to the samples slice.
// Check beep.Streamer interface for further details.
func (m *Synchronized) Stream(samples [][2]float64) (n int, ok bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	n, ok = 0, false
	if m.position >= m.length {
		return n, ok
	}

	n, ok = beep.Mix(m.getStreamers(m.position)...).Stream(samples)
	m.position += n
	if m.position == m.length {
		ok = false
	}
	return n, ok
}

// Err returns an error which occurred during streaming. If no error
// occurred, nil is returned. Check beep.Streamer interface for further
// details
func (m *Synchronized) Err() error {
	return m.err
}

// Add adds audio files into the playlist.
// Calling this method with an already loaded id.
func (m *Synchronized) Add(streamSeeker beep.StreamSeeker, id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, exists := m.streamers[id]; !exists {
		if !m.isInitialized {
			m.length = streamSeeker.Len()
			m.isInitialized = true
		}
		if err := m.validate(streamSeeker); err != nil {
			return fmt.Errorf("cannot load streamer with id %s: %s", id, err)
		}
		m.streamers[id] = &streamSeeker
	}
	return nil
}

// Remove removes audio files from the playlist. It is safe to call
// this method with paths that don't exist in the playlist.
func (m *Synchronized) Remove(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.streamers, id)
}

// validate checks if streamer has the same number of samples as the
// existing streamers. A lock is assumed to be acquired by the caller of
// this method.
func (m *Synchronized) validate(streamSeeker beep.StreamSeeker) error {
	if m.length != streamSeeker.Len() {
		return fmt.Errorf(
			"expected streamer with length %d, got length %d",
			m.length,
			streamSeeker.Len(),
		)
	}
	return nil
}

// getStreamers gets active streamers at the given position. This method calls
// Seek on every streamer. A lock is assumed to be acquired by the caller of this method.
func (m *Synchronized) getStreamers(p int) []beep.Streamer {
	streamers := make([]beep.Streamer, len(m.streamers), len(m.streamers))
	i := 0
	for _, streamer := range m.streamers {
		for id, streamer := range m.streamers {
			if err := (*streamer).Seek(p); err != nil {
				m.err = fmt.Errorf("cannot seek streamer %s to position %d", id, p)
				continue
			}
		}
		streamers[i] = *streamer
		i++
	}
	return streamers
}

func NewMixer() *Synchronized {
	return &Synchronized{
		streamers:     make(map[string]*beep.StreamSeeker),
		position:      0,
		length:        0,
		err:           nil,
		mutex:         sync.RWMutex{},
		isInitialized: false,
	}
}
