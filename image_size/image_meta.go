package imageSize

import (
	"bufio"
	"errors"
	"io"
	"sync"
	"sync/atomic"
)

type Meta struct {
	Format        string
	Width, Height int
}

type format struct {
	magic      string
	decodeMeta func(io.Reader) (*Meta, error)
}

type reader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

var (
	formatsMu     sync.Mutex
	atomicFormats atomic.Value

	ErrFormat = errors.New("unknown image format")
)

func asReader(r io.Reader) reader {
	if rr, ok := r.(reader); ok {
		return rr
	}
	return bufio.NewReader(r)
}

func matchMagic(magic string, b []byte) bool {
	if len(magic) != len(b) {
		return false
	}
	for i, c := range b {
		if magic[i] != c && magic[i] != '?' {
			return false
		}
	}
	return true
}

func RegisterFormat(magic string, decodeMeta func(io.Reader) (*Meta, error)) {
	formatsMu.Lock()
	defer formatsMu.Unlock()

	formats, _ := atomicFormats.Load().([]format)
	atomicFormats.Store(append(formats, format{magic, decodeMeta}))
}

func DecodeMeta(r io.Reader) (*Meta, error) {
	rr := asReader(r)
	formats, _ := atomicFormats.Load().([]format)

	for _, f := range formats {
		b, err := rr.Peek(len(f.magic))
		if err == nil && matchMagic(f.magic, b) {
			return f.decodeMeta(rr)
		}
	}

	return nil, ErrFormat
}
