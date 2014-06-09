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
	w            io.Writer
	platformType PlatformType
	lastByte     byte
}

func NewWriter(w io.Writer, platformType PlatformType) *Writer {
	if platformType != Unix && platformType != Windows {
		if runtime.GOOS == "windows" {
			platformType = Windows
		} else {
			platformType = Unix
		}
	}

	return &Writer{w, platformType, '\000'}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	b := make([]byte, 0)

	var fn func(a byte)

	if w.platformType == Windows {
		fn = func(c byte) {
			if w.lastByte != '\r' && c == '\n' {
				b = append(b, '\r')
				b = append(b, '\n')
			} else {
				b = append(b, c)
			}
		}
	} else {
		fn = func(c byte) {
			if w.lastByte == '\r' && c == '\n' {
				b[len(b)-1] = '\n'
			} else {
				b = append(b, c)
			}
		}
	}

	for _, c := range p {
		fn(c)
		w.lastByte = c
	}

	return w.w.Write(b)
}
