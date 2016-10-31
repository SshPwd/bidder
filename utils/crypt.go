package utils

import (
	"crypto/aes"
	"fmt"
)

var (
	encodeKey = []byte("miloAds123!12345")
)

func Encrpyt(src []byte) (dst []byte) {

	blockSize := ((len(src)-1)/aes.BlockSize + 1) * aes.BlockSize

	if blockSize > len(src) {
		tmp := make([]byte, blockSize)
		copy(tmp, src)
		src = tmp
	}

	cipher, err := aes.NewCipher(encodeKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	dst = make([]byte, blockSize)

	src0, dst0 := src, dst
	for len(src0) > 0 {
		cipher.Encrypt(dst0, src0[:aes.BlockSize])
		dst0, src0 = dst0[aes.BlockSize:], src0[aes.BlockSize:]
	}
	return
}

func Decrpyt(src []byte) (dst []byte) {

	if len(src)%aes.BlockSize != 0 {
		return
	}

	cipher, err := aes.NewCipher(encodeKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	dst = make([]byte, len(src))

	src0, dst0 := src, dst
	for len(src0) > 0 {
		cipher.Decrypt(dst0, src0)
		dst0, src0 = dst0[aes.BlockSize:], src0[aes.BlockSize:]
	}
	return
}
