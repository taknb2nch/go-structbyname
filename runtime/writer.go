package runtime

import (
	"io"
	"runtime"
)

type PlatformType uint32

const (
	AutoDetect PlatformType = iota
	Unix
	Windows
)

type Writer struct {
	w             io.Writer
	lineSeparator []byte
	lastByte      byte
	//platformType  PlatformType
}

func NewWriter(w io.Writer, platformType PlatformType) *Writer {
	var lineSeparator []byte

	switch platformType {
	case AutoDetect:
		if runtime.GOOS == "windows" {
			lineSeparator = []byte{'\r', '\n'}
		} else {
			lineSeparator = []byte{'\n'}
		}
	case Unix:
		lineSeparator = []byte{'\n'}
	case Windows:
		lineSeparator = []byte{'\r', '\n'}
	default:
		lineSeparator = []byte{'\n'}
	}

	return &Writer{w, lineSeparator, '\000'}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	b := make([]byte, 0)
	for _, c := range p {
		if w.lastByte != '\r' && c == '\n' {
			b = append(b, w.lineSeparator...)
		} else {
			b = append(b, c)
		}

		w.lastByte = c
	}

	return w.w.Write(b)
}
