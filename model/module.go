package model

type Module interface {
	Read(request *Request) error
	Write(request *Request) error
	Win(auctionPrice, revenuePrice float64) float64
	Name() string
}
