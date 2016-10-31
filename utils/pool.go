package utils

import (
	"bytes"
	"sync"
)

var poolBuffers sync.Pool = sync.Pool{
	New: func() interface{} {
		buf := new(bytes.Buffer)
		buf.Grow(512)
		return buf
	},
}

func GetBuffer() (buf *bytes.Buffer) {
	return poolBuffers.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
	buf.Truncate(0)
	poolBuffers.Put(buf)
}
