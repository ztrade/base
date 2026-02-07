package common

import (
	"errors"

	"github.com/shopspring/decimal"
	. "github.com/ztrade/trademodel"
)

var (
	ErrNoBalance = errors.New("no balance")
	Zero         = decimal.NewFromInt(0)
)

type VBalance struct {
	total          decimal.Decimal
	prevRoundTotal decimal.Decimal
	position       decimal.Decimal
	feeTotal       decimal.Decimal
	//  开仓的总价值
	longCost  decimal.Decimal
	shortCost decimal.Decimal
	fee       decimal.Decimal
	prevFee   decimal.Decimal
}

func NewVBalance() *VBalance {
	b := new(VBalance)
	b.total = decimal.NewFromFloat(100000)
	b.prevRoundTotal = b.total
	b.fee = decimal.NewFromFloat(0.00075)
	return b
}

func (b *VBalance) Set(total float64) {
	b.total = decimal.NewFromFloat(total)
	b.prevRoundTotal = b.total
}

func (b *VBalance) SetFee(fee float64) {
	b.fee = decimal.NewFromFloat(fee)
}

func (b *VBalance) Pos() (pos float64) {
	pos, _ = b.position.Float64()
	return
}

func (b *VBalance) Get() (total float64) {
	// return b.total + b.costTotal
	total, _ = b.total.Float64()
	return
}

func (b *VBalance) GetFeeTotal() (fee float64) {
	fee, _ = b.feeTotal.Float64()
	return
}

func (b *VBalance) AvgOpenPriceDec() (price decimal.Decimal) {
	switch b.position.Sign() {
	case -1:
		return b.shortCost.Div(b.position.Abs())
	case 0:
		return
	case 1:
		return b.longCost.Div(b.position.Abs())
	}
	return
}

func (b *VBalance) AvgOpenPrice() (price float64) {
	price, _ = b.AvgOpenPriceDec().Float64()
	return
}

func (b *VBalance) AddTrade(tr Trade) (profit, profitRate, onceFee float64, err error) {
	// 新的Trade会改变仓位，所以先记录之前的均价
	prevAvgOpenPrice := b.AvgOpenPriceDec()
	amount := decimal.NewFromFloat(tr.Amount).Abs()
	// 仓位价值
	cost := amount.Mul(decimal.NewFromFloat(tr.Price)).Abs()
	fee := cost.Mul(b.fee)
	onceFee, _ = fee.Float64()
	costAll, _ := cost.Add(fee).Float64()
	if tr.Action.IsOpen() && costAll > b.Get() {
		err = ErrNoBalance
		return
	}
	// close/stop just return if no position
	if b.position.Equal(Zero) && !tr.Action.IsOpen() {
		return
	}
	if tr.Action.IsLong() {
		b.position = b.position.Add(amount)
		b.longCost = b.longCost.Add(cost)
	} else {
		b.position = b.position.Sub(amount)
		b.shortCost = b.shortCost.Add(cost)
	}
	isPositionZero := b.position.Equal(Zero)
	if tr.Action.IsOpen() && !isPositionZero {
		b.total = b.total.Sub(cost).Sub(fee)
	}
	b.feeTotal = b.feeTotal.Add(fee)
	// 注意：部分平仓时 b.total 不反映真实余额，仅在仓位完全归零后通过 shortCost-longCost 重新计算
	// 计算盈利
	if isPositionZero {
		totalFee := fee.Add(b.prevFee)
		prof := b.shortCost.Sub(b.longCost).Sub(totalFee)
		b.total = b.prevRoundTotal.Add(prof)
		profit, _ = prof.Float64()
		profitRate, _ = prof.Div(prevAvgOpenPrice).Float64()
		b.longCost = decimal.NewFromInt(0)
		b.shortCost = decimal.NewFromInt(0)
		b.prevRoundTotal = b.total
		b.prevFee = decimal.Zero
	} else {
		b.prevFee = b.prevFee.Add(fee)
	}
	return
}
