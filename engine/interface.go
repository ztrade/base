package engine

import (
	"github.com/ztrade/base/common"
	"github.com/ztrade/indicator"
	"github.com/ztrade/trademodel"
)

const (
	StatusRunning = 0
	StatusSuccess = 1
	StatusFail    = -1
)

type Engine interface {
	OpenLong(price, amount float64) string
	CloseLong(price, amount float64) string
	OpenShort(price, amount float64) string
	CloseShort(price, amount float64) string
	StopLong(price, amount float64) string
	StopShort(price, amount float64) string
	CancelOrder(string)
	CancelAllOrder()
	DoOrder(typ trademodel.TradeType, price, amount float64) string
	AddIndicator(name string, params ...int) (ind indicator.CommonIndicator)
	Position() (pos, price float64)
	Balance() float64
	Log(v ...interface{})
	Watch(watchType string)
	SendNotify(title, content, contentType string)
	Merge(src, dst string, fn common.CandleFn)
	SetBalance(balance float64)
	UpdateStatus(status int, msg string)
}
