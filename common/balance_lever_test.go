package common

import (
	"testing"

	. "github.com/ztrade/trademodel"
)

func TestCheckLiquidationLong(t *testing.T) {
	lb := NewLeverBalance()
	lb.Set(100)
	lb.SetFee(0.0002)
	lb.SetLever(10)
	_, _, err := lb.AddTrade(Trade{Action: OpenLong, Price: 100, Amount: 9})
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
	_, _, err := lb.AddTrade(Trade{Action: OpenShort, Price: 100, Amount: 9})
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
	_, _, err := lb.AddTrade(Trade{Action: OpenLong, Price: 100, Amount: 10})
	if err == nil {
		t.Fatal("Liq not work")
	}
	t.Log(err.Error())
}
