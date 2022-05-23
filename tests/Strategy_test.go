package tests

import (
	"MoTrade/OKXClient"
	mo_errors "MoTrade/mo-errors"
	"MoTrade/strategies"
	"errors"
	"testing"
	"time"
)

var (
	strategy *strategies.Strategy
)

func PlaceOneStrategy() {
	strategy = strategies.NewStrategy(&client.Trade, 1)
	go strategy.Watching(time.Second / 2)
}

func TestPlaceOneOrder(t *testing.T) {
	PlaceOneStrategy()
	order, err := strategies.NewOrder(&client.Trade, OKXClient.SWAP, OKXClient.ETH_USDT_SWAP, OKXClient.CROSS, OKXClient.SHORT, OKXClient.MARKET, 1, 0)
	if err != nil {
		t.Error(err)
	}
	err = strategy.FillOneOrder(order, time.Second)
	if err != nil {
		switch {
		case errors.Is(err, mo_errors.FullError):
			t.Log("FullError")
			t.Error(err)
		default:
			t.Error(err)
		}
	}
	o, err := strategies.NewOrder(&client.Trade, OKXClient.SWAP, OKXClient.ETH_USDT_SWAP, OKXClient.CROSS, OKXClient.SHORT, OKXClient.MARKET, 1, 0)
	if err != nil {
		t.Error(err)
	}
	err = strategy.FillOneOrder(o, time.Second)
	if err != nil {
		switch {
		case errors.Is(err, mo_errors.FullError):
			t.Log("FullError")
			t.Error(err)
		default:
			t.Error(err)
		}
	}
	select {}
}
