package bidder

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"golang_bidder/config"
	"golang_bidder/logs"
	"golang_bidder/utils"
	"net"
	"sync"
	"time"
)

const (
	C0 = "\x1b[30;1m"
	CR = "\x1b[31;1m"
	CG = "\x1b[32;1m"
	CY = "\x1b[33;1m"
	CB = "\x1b[34;1m"
	CM = "\x1b[35;1m"
	CC = "\x1b[36;1m"
	CW = "\x1b[37;1m"
	CN = "\x1b[0m"
)

func Serve() {

	conf := config.Get()

	listeners := make([]net.Listener, 0, len(conf.Bind))

	for _, bind := range conf.Bind {

		fmt.Println(bind)

		listener, err := reuseport.Listen("tcp4", bind)
		if err != nil {
			logs.Critical(err.Error())
			return
		}

		listeners = append(listeners, listener)
	}

	utils.SetNofile(128000)
	utils.Renice(-20)

	wg := sync.WaitGroup{}

	for _, listener := range listeners {

		wg.Add(1)

		go func(listener net.Listener) {

			serv := fasthttp.Server{
				Name:            "deximedia",
				Handler:         HttpHandle,
				ReadTimeout:     20 * time.Second,
				WriteTimeout:    20 * time.Second,
				WriteBufferSize: 96 * 1024,
				ReadBufferSize:  24 * 1024,
				Concurrency:     42000,
			}

			err := serv.Serve(listener)

			logs.Critical(err.Error())
			wg.Done()

		}(listener)
	}

	fmt.Println("READY")

	time.Sleep(5 * time.Second)

	if prevPid := utils.ReadPid(conf.PidFile); prevPid != 0 {
		fmt.Println("kill", prevPid)
		utils.Kill(prevPid)
	}

	err := utils.WritePid(conf.PidFile)
	if err != nil {
		logs.Critical(err.Error())
	}

	wg.Wait()
}
