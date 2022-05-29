package MABased

import (
	"MoTrade/OKXClient"
	mlog "MoTrade/mo-log"
	"MoTrade/strategies"
	"container/list"
	"time"
)

var log = mlog.Log

type SimpleMABasedStrategy struct {
	strategies.Strategy
	InstType    string
	InstId      string
	TdMode      string
	Size        int
	Bar         string
	Limit       int
	MAs         *list.List // float64
	Prices      *list.List // float64
	TrendUnit   float64
	Observable  float64
	SharpUnit   float64
	Protected   bool
	ProtectLoss float64
	LongPoint   chan struct{}
	ShortPoint  chan struct{}
}

func NewMABasedStrategy(trade *OKXClient.Trade, maxOrder int, instType, instId, tdMode string, size int, bar string, limit int, protected bool, protectLoss float64) *SimpleMABasedStrategy {
	strategy := &SimpleMABasedStrategy{
		Strategy:    *strategies.NewStrategy(trade, maxOrder),
		InstType:    instType,
		InstId:      instId,
		TdMode:      tdMode,
		Size:        size,
		Bar:         bar,
		Limit:       limit,
		MAs:         list.New(),
		Prices:      list.New(),
		TrendUnit:   1,
		Observable:  0.2,
		SharpUnit:   10,
		Protected:   protected,
		ProtectLoss: protectLoss,
		LongPoint:   make(chan struct{}),
		ShortPoint:  make(chan struct{}),
	}
	go strategy.Watching(time.Second / 2)
	return strategy
}

func (strategy *SimpleMABasedStrategy) Run(interval time.Duration) {

	if err := strategy.initLists(); err != nil {
		log.Fatalln("initLists error: %v", err)
		return
	}

	go strategy.PointWatching(interval)

	for {
		select {
		case <-strategy.LongPoint:
			log.Debugln("Long..")
			strategy.KillAllShortOrders(3 * time.Second)
			order, err := strategy.newLongOrder()
			if err != nil {
				log.Println(err)
			}
			if err = strategy.FillOneOrder(order, time.Second/2); err != nil {
				log.Println(err)
			}
		case <-strategy.ShortPoint:
			log.Debugln("Short..")
			strategy.KillAllLongOrders(3 * time.Second)
			order, err := strategy.newShortOrder()
			if err != nil {
				log.Println(err)
			}
			if err = strategy.FillOneOrder(order, time.Second/2); err != nil {
				log.Println(err)
			}
		}
	}
}

func (strategy *SimpleMABasedStrategy) initLists() error {
	strategy.MAs.Init()
	strategy.Prices.Init()
	for i := 1; i < 4; i++ {
		ma, err := strategy.Trade.Market.GetMA(strategy.InstId, strategy.Bar, strategy.Limit, i)
		if err != nil {
			log.Errorln("Get MA Error:", err)
			return err
		}
		strategy.MAs.PushBack(ma)
		value, err := strategy.Trade.Market.GetTickerValue(strategy.InstId)
		if err != nil {
			log.Errorln("Get Ticker Value Error:", err)
			return err
		}
		strategy.Prices.PushBack(value)
	}
	return nil
}

func (strategy *SimpleMABasedStrategy) PointWatching(interval time.Duration) {
	dur, _ := time.ParseDuration(strategy.Bar)
	durTimer := time.NewTimer(dur)
	optTimer := time.NewTimer(interval)
	for {
		select {
		case <-durTimer.C:
			log.Debugln("1 min passed")
			if err := strategy.updateMaAndPrice(); err != nil {
				log.Errorln("updateMaAndPrice Error:", err)
				_ = strategy.initLists()
				durTimer.Reset(0)
				continue
			}
			durTimer.Reset(dur)
		case <-optTimer.C:
			log.Debugln("1 opt")
			res, err := strategy.longOrShort()
			if err != nil {
				log.Errorln("longOrShort Error:", err)
				optTimer.Reset(interval)
				continue
			}
			if res == 1 {
				log.Println("Long Point")
				strategy.LongPoint <- struct{}{}
			} else if res == -1 {
				log.Println("Short Point")
				strategy.ShortPoint <- struct{}{}
			} else {
				log.Debugln("No Point")
			}
			optTimer.Reset(interval)
		}
	}
}

func (strategy *SimpleMABasedStrategy) longOrShort() (int, error) {
	if err := strategy.getMaAndPrice(); err != nil {
		log.Errorln("getMaAndPrice Error:", err)
		return 0, err
	}
	maTrende, priceTrende := 0, 0
	ma0 := strategy.MAs.Front().Value.(float64)
	ma1 := strategy.MAs.Front().Next().Value.(float64)
	ma2 := strategy.MAs.Back().Value.(float64)
	price0 := strategy.Prices.Front().Value.(float64)
	price1 := strategy.Prices.Front().Next().Value.(float64)
	price2 := strategy.Prices.Back().Value.(float64)
	dma := ma2 - ma1
	dprice := price2 - price1

	if ma2-ma0 > strategy.TrendUnit {
		maTrende = 1
	} else if ma0-ma2 > strategy.TrendUnit {
		maTrende = -1
	}
	if price2-price0 > strategy.TrendUnit {
		priceTrende = 1
	} else if price0-price2 > strategy.TrendUnit {
		priceTrende = -1
	}
	log.Debugln("ma0:", ma0, "ma1:", ma1, "ma2:", ma2, "dma:", dma, "maTrende:", maTrende)
	log.Debugln("price0:", price0, "price1:", price1, "price2:", price2, "dprice:", dprice, "priceTrende:", priceTrende)

	if maTrende == 1 && priceTrende == 1 && price1-ma1 < strategy.Observable && price2-ma2 >= strategy.Observable {
		log.Debugln("Long 1: Up Up go")
		return 1, nil
	}

	if maTrende == 1 && priceTrende == -1 && price2-ma2 > strategy.Observable && price2-ma2 < 2*strategy.Observable {
		log.Debugln("Long 2: Up Down go")
		return 1, nil
	}

	if maTrende == 1 && price1-ma1 > strategy.Observable && ma2-price2 > strategy.Observable && price1-price2 > strategy.SharpUnit {
		log.Debugln("Long 3: SharpDown go")
		return 1, nil
	}

	if maTrende != 1 && priceTrende == -1 && ma1-price1 < strategy.Observable && ma2-price2 >= strategy.Observable {
		log.Debugln("Short 1: Down Down run")
		return -1, nil
	}

	if priceTrende == 1 && ma1-price1 > strategy.Observable && ma2-price2 < 2*strategy.Observable && ma2-price2 > strategy.Observable {
		log.Debugln("Short 2: Down Up run")
		return -1, nil
	}

	if priceTrende == 1 && maTrende == 1 && dprice-dma > strategy.SharpUnit {
		log.Debugln("Short 3: SharpUp run")
		return -1, nil
	}
	return 0, nil
}

func (strategy *SimpleMABasedStrategy) newLongOrder() (*strategies.Order, error) {
	order, err := strategies.NewOrder(strategy.Trade, strategy.InstType, strategy.InstId, strategy.TdMode, OKXClient.LONG, OKXClient.MARKET, strategy.Size, 0)
	if strategy.Protected {
		go order.Protect(strategy.ProtectLoss, time.Second/2)
	}
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return order, nil
}

func (strategy *SimpleMABasedStrategy) newShortOrder() (*strategies.Order, error) {
	order, err := strategies.NewOrder(strategy.Trade, strategy.InstType, strategy.InstId, strategy.TdMode, OKXClient.SHORT, OKXClient.MARKET, strategy.Size, 0)
	if strategy.Protected {
		go order.Protect(strategy.ProtectLoss, time.Second/2)
	}
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return order, nil
}

func (strategy *SimpleMABasedStrategy) updateMaAndPrice() error {
	ma, err := strategy.Trade.Market.GetMA(strategy.InstId, strategy.Bar, strategy.Limit, 0)
	if err != nil {
		log.Errorln("Get MA Error:", err)
		return err
	}
	strategy.MAs.PushBack(ma)
	value, err := strategy.Trade.Market.GetTickerValue(strategy.InstId)
	if err != nil {
		log.Errorln("Get Ticker Value Error:", err)
		return err
	}
	strategy.Prices.PushBack(value)
	strategy.MAs.Remove(strategy.MAs.Front())
	strategy.Prices.Remove(strategy.Prices.Front())
	return nil
}

func (strategy *SimpleMABasedStrategy) getMaAndPrice() error {
	ma, err := strategy.Trade.Market.GetMA(strategy.InstId, strategy.Bar, strategy.Limit, 0)
	if err != nil {
		log.Errorln("Get MA Error:", err)
		return err
	}
	price, err := strategy.Trade.Market.GetTickerValue(strategy.InstId)
	if err != nil {
		log.Errorln("Get Ticker Value Error:", err)
		return err
	}
	strategy.MAs.Back().Value = ma
	strategy.Prices.Back().Value = price
	return nil
}
