package utils

import (
	"syscall"
)

func SetNofile(n uint64) {

	syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{Cur: n, Max: n})
}
