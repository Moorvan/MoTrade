package strategies

import (
	"MoTrade/OKXClient"
	mlog "MoTrade/mo-log"
	"time"
)

var (
	log = mlog.Log
)

type Strategy struct {
	Trade      *OKXClient.Trade
	Orders     []*Order
	MaxOrder   int
	InstId     string
	InstType   string
	TdMode     string
	Profit     float64
	ApprProfit float64
}

// TODO: test

func (strategy *Strategy) KillAllOrders() {
	for _, order := range strategy.Orders {
		if err := order.CleanOrder(OKXClient.MARKET, 0, 5*time.Second); err != nil {
			log.Alarm("KillOrders Failed", err)
		} else {
			if err = log.WriteLog("killedOrders", order); err != nil {
				log.Errorln("WriteLog Failed", err)
			}
		}
		if order.IsFinished {
			strategy.Profit += order.Profit
		}
	}
}

func (strategy *Strategy) FillOneOrder(posSide string, size int, timeout time.Duration) (bool, error) {
	if len(strategy.Orders) >= strategy.MaxOrder {
		return false, nil
	}
	order, err := NewOrder(strategy.Trade, strategy.InstType, strategy.InstId, OKXClient.CROSS, posSide, OKXClient.MARKET, size, 0, timeout)
	if err != nil {
		return false, err
	}
	strategy.Orders = append(strategy.Orders, order)
	return true, nil
}

func (strategy *Strategy) FillOneOrderWithLimit(posSide string, size int, timeout time.Duration, limit float64) (bool, error) {
	if len(strategy.Orders) >= strategy.MaxOrder {
		return false, nil
	}
	order, err := NewOrder(strategy.Trade, strategy.InstType, strategy.InstId, OKXClient.CROSS, posSide, OKXClient.LIMIT, size, limit, timeout)
	if err != nil {
		return false, err
	}
	strategy.Orders = append(strategy.Orders, order)
	return true, nil
}
