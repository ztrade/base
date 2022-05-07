package engine

import (
	"github.com/ztrade/base/common"
	"github.com/ztrade/indicator"
)

const (
	StatusRunning = 0
	StatusSuccess = 1
	StatusFail    = -1
)

type Engine interface {
	OpenLong(price, amount float64)
	CloseLong(price, amount float64)
	OpenShort(price, amount float64)
	CloseShort(price, amount float64)
	StopLong(price, amount float64)
	StopShort(price, amount float64)
	CancelAllOrder()
	AddIndicator(name string, params ...int) (ind indicator.CommonIndicator)
	Position() (pos, price float64)
	Balance() float64
	Log(v ...interface{})
	Watch(watchType string)
	SendNotify(content, contentType string)
	Merge(src, dst string, fn common.CandleFn)
	SetBalance(balance float64)
	UpdateStatus(status int, msg string)
}
