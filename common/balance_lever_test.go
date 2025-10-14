package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	. "github.com/ztrade/trademodel"
)

func TestCheckLiquidationLong(t *testing.T) {
	lb := NewLeverBalance()
	lb.Set(100)
	lb.SetFee(0.0002)
	lb.SetLever(10)
	_, _, _, err := lb.AddTrade(Trade{Action: OpenLong, Price: 100, Amount: 9})
	if err != nil {
		t.Fatal("Liq lever AddTrade failed:" + err.Error())
	}
	liqPrice, isLiq := lb.CheckLiquidation(90.1)
	if isLiq {
		t.Fatal("Liq cal too large")
	}
	t.Log(liqPrice, isLiq)
	liqPrice, isLiq = lb.CheckLiquidation(90)
	if !isLiq {
		t.Fatal("Liq cal too small")
	}
	t.Log(liqPrice, isLiq)
}

func TestCheckLiquidationShort(t *testing.T) {
	lb := NewLeverBalance()
	lb.Set(100)
	lb.SetFee(0.0002)
	lb.SetLever(10)
	_, _, _, err := lb.AddTrade(Trade{Action: OpenShort, Price: 100, Amount: 9})
	if err != nil {
		t.Fatal("Liq lever AddTrade failed:" + err.Error())
	}
	liqPrice, isLiq := lb.CheckLiquidation(109)
	if isLiq {
		t.Fatal("Liq cal too small")
	}
	t.Log(liqPrice, isLiq)
	liqPrice, isLiq = lb.CheckLiquidation(109.99)
	if !isLiq {
		t.Fatal("Liq cal too large")
	}
	t.Log(liqPrice, isLiq)
}

func TestCheckLeverBalance(t *testing.T) {
	lb := NewLeverBalance()
	lb.Set(100)
	lb.SetFee(0.0002)
	lb.SetLever(10)
	_, _, _, err := lb.AddTrade(Trade{Action: OpenLong, Price: 100, Amount: 10})
	if err == nil {
		t.Fatal("Liq not work")
	}
	t.Log(err.Error())
}

func TestLeverLong(t *testing.T) {
	tm := time.Now()
	openTrade := Trade{
		ID:     "1",
		Action: OpenLong,
		Time:   tm,
		Price:  100,
		Amount: 1,
	}
	closeTrade := Trade{
		ID:     "2",
		Action: CloseLong,
		Time:   tm.Add(time.Second),
		Price:  110,
		Amount: 1,
	}
	stopTrade := Trade{
		ID:     "3",
		Action: StopLong,
		Time:   tm.Add(time.Second * 2),
		Price:  95,
		Amount: 1,
	}

	t.Log("---- test lever long ----")
	b := NewLeverBalance()
	b.Set(100)
	b.SetFee(0.0002)
	b.SetLever(10)
	_, _, onceFee, err := b.AddTrade(openTrade)
	assert.NoError(t, err)
	assert.Equal(t, b.Get(), 89.98)
	assert.Equal(t, onceFee, 0.02)
	t.Logf("balance after open trade: %f, fee: %f", b.Get(), onceFee)
	profit, _, onceFee, err := b.AddTrade(closeTrade)
	assert.NoError(t, err)
	t.Log("profit:", profit, onceFee)
	t.Logf("balance after close trade: %f, fee: %f", b.Get(), onceFee)
	fee := calFee(b.vBalance.fee, openTrade, closeTrade)
	t.Log("fee:", fee)
	if b.Get() != 110-fee {
		t.Fatal("balance close error:", b.Get(), 1010-fee)
	}
	assert.Equal(t, profit, 10.0-fee)

	t.Log("---- test stop long ----")
	b = NewLeverBalance()
	b.Set(100)
	b.SetFee(0.0002)
	b.SetLever(10)
	_, _, onceFee, err = b.AddTrade(openTrade)
	assert.NoError(t, err)
	t.Logf("balance after open trade: %f, fee: %f", b.Get(), onceFee)

	profit, _, _, err = b.AddTrade(stopTrade)
	assert.NoError(t, err)
	t.Logf("balance after stop trade: %f, fee: %f", b.Get(), onceFee)
	t.Log("profit:", profit, onceFee)
	fee = calFee(b.vBalance.fee, openTrade, stopTrade)
	t.Log("fee:", fee)
	if b.Get() != 95-fee {
		t.Fatal("balance stop error:", b.Get())
	}
	assert.Equal(t, profit, -5.0-fee)

	t.Log("---- test liquidation when stop long  ----")
	b = NewLeverBalance()
	b.Set(100)
	b.SetFee(0.0002)
	b.SetLever(10)
	_, _, openFee, err := b.AddTrade(openTrade)
	assert.NoError(t, err)
	t.Logf("balance after open trade: %f, fee: %f", b.Get(), onceFee)
	stopTrade.Price = 80
	profit, _, onceFee, err = b.AddTrade(stopTrade)
	assert.NoError(t, err)
	t.Logf("balance after stop trade: %f, fee: %f", b.Get(), onceFee)
	t.Log("profit:", profit, onceFee)
	// when liquidation: balance = startBalance - openFee
	assert.Equal(t, b.Get(), 100-10-openFee)
}

func TestLeverShort(t *testing.T) {
	tm := time.Now()
	openTrade := Trade{
		ID:     "1",
		Action: OpenShort,
		Time:   tm,
		Price:  100,
		Amount: 1,
	}
	closeTrade := Trade{
		ID:     "2",
		Action: CloseShort,
		Time:   tm.Add(time.Second),
		Price:  90,
		Amount: 1,
	}
	stopTrade := Trade{
		ID:     "3",
		Action: StopShort,
		Time:   tm.Add(time.Second * 2),
		Price:  105,
		Amount: 1,
	}

	t.Log("---- test lever short ----")
	b := NewLeverBalance()
	b.Set(100)
	b.SetFee(0.0002)
	b.SetLever(10)
	_, _, onceFee, err := b.AddTrade(openTrade)
	assert.NoError(t, err)
	assert.Equal(t, b.Get(), 89.98)
	assert.Equal(t, onceFee, 0.02)
	t.Logf("balance after open trade: %f, fee: %f", b.Get(), onceFee)
	profit, _, onceFee, err := b.AddTrade(closeTrade)
	assert.NoError(t, err)
	t.Log("profit:", profit, onceFee)
	t.Logf("balance after close trade: %f, fee: %f", b.Get(), onceFee)
	fee := calFee(b.vBalance.fee, openTrade, closeTrade)
	t.Log("fee:", fee)
	if b.Get() != 110-fee {
		t.Fatal("balance close error:", b.Get(), 1010-fee)
	}
	assert.Equal(t, profit, 10.0-fee)

	t.Log("---- test stop short ----")
	b = NewLeverBalance()
	b.Set(100)
	b.SetFee(0.0002)
	b.SetLever(10)
	_, _, onceFee, err = b.AddTrade(openTrade)
	assert.NoError(t, err)
	t.Logf("balance after open trade: %f, fee: %f", b.Get(), onceFee)

	profit, _, _, err = b.AddTrade(stopTrade)
	assert.NoError(t, err)
	t.Logf("balance after stop trade: %f, fee: %f", b.Get(), onceFee)
	t.Log("profit:", profit, onceFee)
	fee = calFee(b.vBalance.fee, openTrade, stopTrade)
	t.Log("fee:", fee)
	if b.Get() != 95-fee {
		t.Fatal("balance stop error:", b.Get())
	}
	assert.Equal(t, profit, -5.0-fee)

	t.Log("---- test liquidation when stop short  ----")
	b = NewLeverBalance()
	b.Set(100)
	b.SetFee(0.0002)
	b.SetLever(10)
	_, _, openFee, err := b.AddTrade(openTrade)
	assert.NoError(t, err)
	t.Logf("balance after open trade: %f, fee: %f", b.Get(), openFee)
	stopTrade.Price = 120
	profit, _, onceFee, err = b.AddTrade(stopTrade)
	assert.NoError(t, err)
	t.Logf("balance after stop trade: %f, fee: %f", b.Get(), onceFee)
	t.Log("profit:", profit, onceFee)
	// when liquidation: balance = startBalance - openFee
	// Floating point number calculation accuracy loss
	diff := b.Get() - (100 - 10 - openFee)
	assert.InDelta(t, diff, 0, 0.0001)
}
