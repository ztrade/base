package common

import (
	"github.com/shopspring/decimal"
	. "github.com/ztrade/trademodel"
)

var (
	numOne = decimal.NewFromInt(1)
)

type LeverBalance struct {
	vBalance *VBalance
	total    decimal.Decimal

	//  开仓的总价值
	lever decimal.Decimal
}

func NewLeverBalance() *LeverBalance {
	lb := new(LeverBalance)
	lb.vBalance = NewVBalance()
	lb.lever = decimal.NewFromFloat(1)
	return lb
}

func (b *LeverBalance) Set(total float64) {
	b.total = decimal.NewFromFloat(total)
	vTotal, _ := b.total.Mul(b.lever).Float64()
	b.vBalance.Set(vTotal)
}

func (b *LeverBalance) SetFee(fee float64) {
	b.vBalance.SetFee(fee)
}

func (b *LeverBalance) SetLever(lever float64) {
	b.lever = decimal.NewFromFloat(lever)
	vTotal, _ := b.total.Mul(b.lever).Float64()
	b.vBalance.Set(vTotal)
}

func (b *LeverBalance) Pos() (pos float64) {
	return b.vBalance.Pos()
}

// func (b *LeverBalance) LiquidationPrice() (price float64, valid bool) {
// 	pos, _ := b.position.Float64()
// 	if pos == 0 {
// 		return
// 	}
// 	valid = true
// 	if pos > 0 {
// 		price, _ = b.openPrice.Sub(b.openPrice.Div(b.lever)).Float64()
// 	} else {
// 		price, _ = b.openPrice.Add(b.openPrice.Div(b.lever)).Float64()
// 	}
// 	return
// }

func (b *LeverBalance) CheckLiquidation(price float64) (liqPrice float64, isLiq bool) {
	openPrice := decimal.NewFromFloat(b.vBalance.AvgOpenPrice())
	fee := b.vBalance.fee
	switch b.vBalance.position.Sign() {
	// <0
	case -1:
		// liqPrice + liqPrice * fee = openPrice + openPrice/lever
		// liqPrice *(1 + fee) = openPrice * (1 + 1/lever)
		// liqPrice = (openPrice * (1 + 1/lever))/(1-fee)
		liqPrice, _ = openPrice.Add(openPrice.Div(b.lever)).Div(numOne.Add(fee)).Float64()
		if price >= liqPrice {
			isLiq = true
		}
	// =0
	case 0:
		return
	// >0
	case 1:
		// liqPrice - liqPrice * fee = openPrice - openPrice/lever
		// (1-fee) * liqPrice = openPrice * (1 - 1/lever)
		// liqPrice = (openPrice * (1 - 1/lever))/(1-fee)
		liqPrice, _ = openPrice.Sub(openPrice.Div(b.lever)).Div(numOne.Sub(fee)).Float64()
		if price <= liqPrice {
			isLiq = true
		}
	}
	return
}

func (b *LeverBalance) Get() (total float64) {
	total, _ = b.total.Float64()
	return
}

func (b *LeverBalance) GetFeeTotal() float64 {
	return b.vBalance.GetFeeTotal()
}

func (b *LeverBalance) AddTrade(tr Trade) (profit, profitRate, onceFee float64, err error) {
	if tr.Action.IsOpen() {
		// check balance enough when open order
		amount := decimal.NewFromFloat(tr.Amount).Abs()
		cost := amount.Mul(decimal.NewFromFloat(tr.Price)).Abs()
		fee := cost.Mul(b.vBalance.fee)
		onceCost := cost.Div(b.lever).Add(fee)
		if b.total.LessThan(onceCost) {
			err = ErrNoBalance
			return
		}
	} else {
		liqPrice, isLiq := b.CheckLiquidation(tr.Price)
		if isLiq {
			tr.Price = liqPrice
		}
	}
	profit, profitRate, onceFee, err = b.vBalance.AddTrade(tr)
	if err != nil {
		return
	}
	b.total = b.total.Add(decimal.NewFromFloat(profit)).Sub(decimal.NewFromFloat(onceFee))
	return
}
