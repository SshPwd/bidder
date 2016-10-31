package utils

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
	"strings"
)

const TRIM = "\t\r\n "

func ToUint32(val string) uint32 {
	u64, _ := strconv.ParseUint(val, 10, 32)
	return uint32(u64)
}

func ToUint16(val string) uint16 {
	u64, _ := strconv.ParseUint(val, 10, 16)
	return uint16(u64)
}

func ToUint8(val string) uint8 {
	u64, _ := strconv.ParseUint(val, 10, 8)
	return uint8(u64)
}

func ToUint(val string) uint {
	u64, _ := strconv.ParseUint(val, 10, 32)
	return uint(u64)
}

func ToInt(val string) int {
	i64, _ := strconv.ParseInt(val, 10, 32)
	return int(i64)
}

func ToUint64(val string) uint64 {
	u64, _ := strconv.ParseUint(val, 10, 64)
	return u64
}

func ToInt64(val string) int64 {
	i64, _ := strconv.ParseInt(val, 10, 64)
	return i64
}

func ToFloat32(val string) float32 {
	f64, _ := strconv.ParseFloat(val, 32)
	return float32(f64)
}

func ToFloat64(val string) float64 {
	f64, _ := strconv.ParseFloat(val, 64)
	return f64
}

func IpToUint32(ipStr string) (ip uint32) {
	if ipByte := net.ParseIP(ipStr).To4(); ipByte != nil {
		ip4 := []byte(ipByte)
		binary.Read(bytes.NewReader(ip4), binary.BigEndian, &ip)
	}
	return
}

func Uint32ToIp(ip uint32) (ipStr string) {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ip)
	ipStr = net.IPv4(ipByte[0], ipByte[1], ipByte[2], ipByte[3]).String()
	return
}

func ToSliceInt(strArray []string) (tmp []int) {
	tmp = make([]int, 0, len(strArray))
	for _, s := range strArray {
		s = strings.Trim(s, TRIM)
		if val, err := strconv.ParseInt(s, 10, 32); err == nil {
			tmp = append(tmp, int(val))
		}
	}
	return
}

func ToSliceInt64(strArray []string) (tmp []int64) {
	tmp = make([]int64, 0, len(strArray))
	for _, s := range strArray {
		s = strings.Trim(s, TRIM)
		if val, err := strconv.ParseInt(s, 10, 64); err == nil {
			tmp = append(tmp, val)
		}
	}
	return
}
