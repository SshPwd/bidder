package utils

const maxIntChars = 18

func ParseInt(buf []byte) int {
	v, n := parseIntBuf(buf)
	if n != len(buf) {
		return 0
	}
	return v
}

func parseIntBuf(b []byte) (int, int) {
	n := len(b)
	if n == 0 {
		return 0, 0
	}
	v := 0
	for i := 0; i < n; i++ {
		c := b[i]
		k := c - '0'
		if k > 9 {
			if i == 0 {
				return -1, i
			}
			return v, i
		}
		if i >= maxIntChars {
			return 0, i
		}
		v = 10*v + int(k)
	}
	return v, n
}
