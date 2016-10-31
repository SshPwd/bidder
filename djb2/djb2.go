package djb2

import (
	"hash"
)

type djb2StringHash32 uint32

func NewDjb32() hash.Hash32 {

	sh := djb2StringHash32(0)
	sh.Reset()

	return &sh
}

func (sh *djb2StringHash32) Size() int {

	return 4
}

func (sh *djb2StringHash32) BlockSize() int {

	return 1
}

func (sh *djb2StringHash32) Sum32() uint32 {

	return uint32(*sh)
}

func (sh *djb2StringHash32) Reset() {

	*sh = djb2StringHash32(5381)
}

func (sh *djb2StringHash32) Sum(b []byte) []byte {

	p := make([]byte, 4)

	p[0] = byte(*sh >> 24)
	p[1] = byte(*sh >> 16)
	p[2] = byte(*sh >> 8)
	p[3] = byte(*sh)

	if b == nil {
		return p
	}

	return append(b, p...)
}

func (sh *djb2StringHash32) Write(b []byte) (int, error) {

	h := uint32(*sh)

	for _, c := range b {
		h = 33*h + uint32(c)
	}

	*sh = djb2StringHash32(h)

	return len(b), nil
}

func Sum32(str string) uint32 {

	sh := NewDjb32()
	sh.Write([]byte(str))

	return sh.Sum32()
}

func Sum32Bytes(data []byte) uint32 {

	sh := NewDjb32()
	sh.Write(data)

	return sh.Sum32()
}
