package browscap

import (
	"sync"
	"unicode/utf8"
)

var (
	bytesPool = &sync.Pool{}
	minCap    = 128
)

func getBytes(size int) []byte {
	if b := bytesPool.Get(); b != nil {
		bs := b.([]byte)
		if cap(bs) >= size {
			return bs[:size]
		}
	}
	c := size
	if c < minCap {
		c = minCap
	}
	return make([]byte, size, c)
}

func mapToBytes(mapping func(rune) rune, s string) []byte {
	maxbytes := len(s)
	nbytes := 0
	var b []byte
	for i, c := range s {
		r := mapping(c)
		if b == nil {
			if r == c {
				continue
			}
			b = getBytes(maxbytes)
			nbytes = copy(b, s[:i])
		}
		if r >= 0 {
			wid := 1
			if r >= utf8.RuneSelf {
				wid = utf8.RuneLen(r)
			}
			if nbytes+wid > maxbytes {
				maxbytes = maxbytes*2 + utf8.UTFMax
				nb := getBytes(maxbytes)
				copy(nb, b[0:nbytes])
				b = nb
			}
			nbytes += utf8.EncodeRune(b[nbytes:maxbytes], r)
		}
	}
	if b == nil {
		b = getBytes(maxbytes)
		copy(b, s)
		return b
	}
	return b[0:nbytes]
}

func mapBytes(mapping func(r rune) rune, s []byte) []byte {
	maxbytes := len(s)
	nbytes := 0
	b := getBytes(maxbytes)
	for i := 0; i < len(s); {
		wid := 1
		r := rune(s[i])
		if r >= utf8.RuneSelf {
			r, wid = utf8.DecodeRune(s[i:])
		}
		r = mapping(r)
		if r >= 0 {
			rl := utf8.RuneLen(r)
			if rl < 0 {
				rl = len(string(utf8.RuneError))
			}
			if nbytes+rl > maxbytes {
				maxbytes = maxbytes*2 + utf8.UTFMax
				nb := getBytes(maxbytes)
				copy(nb, b[0:nbytes])
				b = nb
			}
			nbytes += utf8.EncodeRune(b[nbytes:maxbytes], r)
		}
		i += wid
	}
	return b[0:nbytes]
}
