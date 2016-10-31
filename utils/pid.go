package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func ReadPid(pidfile string) int {

	data, err := ioutil.ReadFile(pidfile)
	if err != nil {
		return 0
	}

	v, _ := strconv.ParseInt(string(data), 10, 32)
	return int(v)
}

func WritePid(pidfile string) error {

	pid := os.Getpid()

	file, err := os.OpenFile(pidfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%d", pid))
	return err
}
