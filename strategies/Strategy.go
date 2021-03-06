package strategies

import (
	"MoTrade/OKXClient"
	mo_errors "MoTrade/mo-errors"
	mlog "MoTrade/mo-log"
	"github.com/thoas/go-funk"
	"sync"
	"time"
)

var (
	log    = mlog.Log
	record = mlog.PersistentRecord
)

type Strategy struct {
	Trade      *OKXClient.Trade
	Orders     []*Order
	MaxOrder   int
	IsStart    chan struct{}
	IsWatching bool
	Profit     float64
	ApprProfit float64
	Side       string
}

func NewStrategy(trade *OKXClient.Trade, maxOrder int) *Strategy {
	return &Strategy{
		Trade:      trade,
		MaxOrder:   maxOrder,
		IsStart:    make(chan struct{}),
		IsWatching: false,
		Profit:     0,
		ApprProfit: 0,
	}
}

func (strategy *Strategy) KillAllOrders(timeout time.Duration) {
	for _, order := range strategy.Orders {
		if err := order.CleanOrder(OKXClient.MARKET, 0, timeout); err != nil {
			log.Alarm("KillOrders Failed", err)
		} else {
			record.PrintStruct(order)
		}
		if !order.IsFinished {
			log.Fatalln("Can't reach here")
		}
	}
}

func (strategy *Strategy) KillAllLongOrders(timeout time.Duration) {
	for _, order := range strategy.Orders {
		if order.PosSide == OKXClient.LONG {
			if err := order.CleanOrder(OKXClient.MARKET, 0, timeout); err != nil {
				log.Alarm("KillOrders Failed", err)
			} else {
				record.PrintStruct(order)
			}
			if !order.IsFinished {
				log.Fatalln("Can't reach here")
			}
		}
	}
}

func (strategy *Strategy) KillAllShortOrders(timeout time.Duration) {
	for _, order := range strategy.Orders {
		if order.PosSide == OKXClient.SHORT {
			if err := order.CleanOrder(OKXClient.MARKET, 0, timeout); err != nil {
				log.Alarm("KillOrders Failed", err)
			} else {
				record.PrintStruct(order)
			}
			if order.IsFinished {
				strategy.Profit += order.Profit
			} else {
				log.Fatalln("Can't reach here")
			}
		}
	}
}

var once sync.Once

func (strategy *Strategy) FillOneOrder(order *Order, timeout time.Duration) error {
	if len(strategy.Orders) >= strategy.MaxOrder {
		return mo_errors.FullError
	}
	if err := order.Start(timeout); err != nil {
		return err
	}
	strategy.Orders = append(strategy.Orders, order)
	once.Do(func() {
		strategy.IsStart <- struct{}{}
	})
	return nil
}

func (strategy *Strategy) Watching(interval time.Duration) {
	strategy.IsWatching = true
	select {
	case <-strategy.IsStart:
	}
	for {
		t := time.NewTimer(interval)
		var sum float64 = 0
		for _, order := range strategy.Orders {
			if !order.IsWatching {
				go order.Watching(interval)
			}
			sum += order.ApprProfit
		}
		strategy.ApprProfit = sum
		log.Debugln("SumApprProfit", strategy.ApprProfit, "SumProfit", strategy.Profit, "Order Count", len(strategy.Orders), "Side", strategy.Side)
		strategy.Orders = funk.Filter(strategy.Orders, func(order *Order) bool {
			if order.IsFinished {
				strategy.Profit += order.Profit
			}
			return !order.IsFinished
		}).([]*Order)
		select {
		case <-t.C:
			continue
		}
	}
}
