package utils

import (
	"bytes"
	"errors"
	"github.com/valyala/fasthttp"
	"time"
)

var (
	getClient = fasthttp.Client{
		MaxConnsPerHost: 5120,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
	}

	strGet        = []byte("GET")
	errConnection = errors.New("errConnection")
	errStatusCode = errors.New("errStatusCode")
)

func HttpGet(url string, timeout time.Duration) (buf *bytes.Buffer, resErr error) {

	res := fasthttp.AcquireResponse()
	req := fasthttp.AcquireRequest()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(url)
	req.Header.SetMethodBytes(strGet)

	err := getClient.DoTimeout(req, res, timeout)
	if err == nil {

		if res.StatusCode() != 200 {
			return nil, errStatusCode
		}

		buf := GetBuffer()
		buf.Write(res.Body())
		return buf, nil
	}

	return nil, errConnection
}
