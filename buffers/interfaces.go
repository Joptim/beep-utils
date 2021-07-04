package buffers

import "github.com/faiface/beep"

type IBuffers interface {
	Load(path string) error
	Release(path string)
	ReleaseAll()
	GetStreamSeeker(path string) (beep.StreamSeeker, error)
	GetFormat() (beep.Format, error)
}
