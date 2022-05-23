package MABased

import (
	"MoTrade/OKXClient"
	mlog "MoTrade/mo-log"
	"MoTrade/strategies"
	"time"
)

var log = mlog.Log

type SimpleMABasedStrategy struct {
	strategies.Strategy
	InstType string
	InstId   string
	TdMode   string
	Size     int
}

func NewMABasedStrategy(trade *OKXClient.Trade, maxOrder int, instType, instId, tdMode string, size int) *SimpleMABasedStrategy {
	strategy := &SimpleMABasedStrategy{
		Strategy: *strategies.NewStrategy(trade, maxOrder),
		InstType: instType,
		InstId:   instId,
		TdMode:   tdMode,
		Size:     size,
	}
	go strategy.Watching(time.Second / 2)
	return strategy
}

func (strategy *SimpleMABasedStrategy) Run() {
	order, err := strategy.newOrder()
	if err != nil {
		log.Fatalln(err)
		return
	}
	if err = strategy.FillOneOrder(order, time.Second/2); err != nil {
		log.Fatalln(err)
		return
	}

	select {}
}

func (strategy *SimpleMABasedStrategy) newOrder() (*strategies.Order, error) {
	order, err := strategies.NewOrder(strategy.Trade, strategy.InstType, strategy.InstId, strategy.TdMode, OKXClient.LONG, OKXClient.MARKET, strategy.Size, 0)
	go order.Protect(0.2, time.Second/2)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return order, nil
}
