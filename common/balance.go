package common

import (
	"errors"

	"github.com/shopspring/decimal"
	. "github.com/ztrade/trademodel"
)

var (
	ErrNoBalance = errors.New("no balance")
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

func (b *VBalance) AvgOpenPrice() (price float64) {
	switch b.position.Sign() {
	case -1:
		price, _ = b.shortCost.Div(b.position.Abs()).Float64()
	case 0:
		return
	case 1:
		price, _ = b.longCost.Div(b.position.Abs()).Float64()
	}

	return
}

func (b *VBalance) AddTrade(tr Trade) (profit, onceFee float64, err error) {
	amount := decimal.NewFromFloat(tr.Amount).Abs()
	// 仓位价值
	cost := amount.Mul(decimal.NewFromFloat(tr.Price)).Abs()
	fee := cost.Mul(b.fee)
	onceFee, _ = fee.Float64()
	costAll, _ := cost.Add(fee).Float64()
	if tr.Action.IsOpen() && costAll >= b.Get() {
		err = ErrNoBalance
		return
	}
	if tr.Action.IsLong() {
		b.position = b.position.Add(amount)
		b.longCost = b.longCost.Add(cost)
	} else {
		b.position = b.position.Sub(amount)
		b.shortCost = b.shortCost.Add(cost)
	}
	isPositionZero := b.position.Equal(decimal.NewFromInt(0))
	if tr.Action.IsOpen() && !isPositionZero {
		b.total = b.total.Sub(cost).Sub(fee)
	}
	b.feeTotal = b.feeTotal.Add(fee)
	// 计算盈利
	if isPositionZero {
		totalFee := fee.Add(b.prevFee)
		prof := b.shortCost.Sub(b.longCost).Sub(totalFee)
		b.total = b.prevRoundTotal.Add(prof)
		profit, _ = prof.Float64()
		b.longCost = decimal.NewFromInt(0)
		b.shortCost = decimal.NewFromInt(0)
		b.prevRoundTotal = b.total
		b.prevFee = decimal.Zero
	} else {
		b.prevFee = b.prevFee.Add(fee)
	}
	return
}
