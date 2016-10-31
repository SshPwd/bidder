package utils

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"time"
)

// 012345678901234567890123456789012345
// 143285c9-9608-0153-a6e2-4e89cc5ac070

func GUID() (guid [36]byte) {

	unixnano := time.Now().UnixNano()
	dstBytes, hexBytes := [8]byte{}, [16]byte{}

	guid[8], guid[13] = '-', '-'
	guid[18], guid[23] = '-', '-'

	binary.BigEndian.PutUint64(dstBytes[:], uint64(unixnano))
	hex.Encode(hexBytes[:], dstBytes[:])

	copy(guid[:8], hexBytes[:8])
	copy(guid[9:13], hexBytes[8:12])
	copy(guid[14:19], hexBytes[12:16])

	rand.Read(dstBytes[:])
	hex.Encode(hexBytes[:], dstBytes[:])

	copy(guid[19:23], hexBytes[:4])
	copy(guid[24:36], hexBytes[4:])
	return
}
