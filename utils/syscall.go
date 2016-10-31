package utils

import (
	"os"
	"syscall"
)

func Renice(priority int) {

	syscall.Setpriority(syscall.PRIO_PROCESS, os.Getpid(), priority)
}

func Kill(pid int) {

	syscall.Kill(pid, syscall.SIGKILL)
}
