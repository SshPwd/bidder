package utils

import (
	"bytes"
	"encoding/base64"
)

func Encode(src []byte) (dst *bytes.Buffer) {

	dst = GetBuffer()
	encoder := base64.NewEncoder(base64.RawURLEncoding, dst)
	encoder.Write(src)
	encoder.Close()
	return
}

func Decode(src []byte) (dst *bytes.Buffer) {

	dst = GetBuffer()
	decoder := base64.NewDecoder(base64.RawURLEncoding, bytes.NewReader(src))
	dst.ReadFrom(decoder)
	return
}
