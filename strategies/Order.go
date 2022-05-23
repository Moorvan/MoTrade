package strategies

import (
	"MoTrade/OKXClient"
	mo_errors "MoTrade/mo-errors"
	"time"
)

type Order struct {
	Trade        *OKXClient.Trade
	InstId       string
	InstType     string
	TdMode       string
	OrdType      string
	PosSide      string
	OpenOrdId    string
	CloseOrdId   string
	PriceIn      float64
	PriceOut     float64
	Size         int
	UnitSize     float64
	IsStart      bool
	IsFinished   bool
	Profit       float64
	ApprProfit   float64
	IsWatching   bool
	IsProfitable bool
	IsProtect    bool
}

func NewOrder(trade *OKXClient.Trade, instType, instId, tdMode, posSide, ordType string, size int, px float64, timeout time.Duration) (*Order, error) {
	var side string
	if posSide == OKXClient.LONG {
		side = OKXClient.BUY
	} else {
		side = OKXClient.SELL
	}
	ordId, err := trade.Market.PlaceOrder(instId, tdMode, side, posSide, ordType, size, px)
	if err != nil {
		log.Errorln("PlaceOrder error:", err)
		return nil, err
	}
	unitSize, err := trade.Market.GetTickerUnitSize(instType, instId)
	if err != nil {
		log.Errorln("GetTickerUnitSize error:", err)
		return nil, err
	}
	order := &Order{
		Trade:     trade,
		InstId:    instId,
		TdMode:    tdMode,
		OrdType:   ordType,
		PosSide:   posSide,
		OpenOrdId: ordId,
		UnitSize:  unitSize,
		Profit:    0,
	}
	go order.watchOpenOrder(time.After(timeout))
	return order, nil
}

func (order *Order) watchOpenOrder(wait <-chan time.Time) {
	for {
		select {
		case <-wait:
			if err := order.CancelOrder(); err != nil {
				log.Alarm("CancelOrder error:", err)
				return
			}
			return
		default:
		}
		info, err := order.Trade.Market.GetOrderInfo(order.InstId, order.OpenOrdId)
		if err != nil {
			log.Println("GetOrderInfo error:", err)
			continue
		}
		if info.State == OKXClient.FILLED {
			order.PriceIn = info.AvgPx
			order.Profit += info.Fee
			order.Size = info.Size
			order.IsStart = true
			return
		}
	}
}

func (order *Order) CancelOrder() error {
	if err := order.Trade.Market.CancelOrder(order.InstId, order.OpenOrdId); err != nil {
		log.Println("CancelOrder error:", err)
		return err
	}
	info, err := order.Trade.Market.GetOrderInfo(order.InstId, order.OpenOrdId)
	if err != nil {
		log.Println("GetOrderInfo error:", err)
		return err
	}
	order.Size = info.Size
	order.PriceIn = info.AvgPx
	order.Profit += info.Fee
	if order.Size == 0 {
		order.IsFinished = true
	}
	order.IsStart = true
	return nil
}

func (order *Order) CleanOrder(ordType string, px float64, timeout time.Duration) error {
	var side string
	if order.IsFinished {
		return nil
	}
	if err := order.waitForOrderStart(timeout, time.Second/3); err != nil {
		return err
	}
	if order.PosSide == OKXClient.LONG {
		side = OKXClient.SELL
	} else {
		side = OKXClient.BUY
	}
	orderId, err := order.Trade.Market.PlaceOrder(order.InstId, order.TdMode, side, order.PosSide, ordType, order.Size, px)
	if err != nil {
		log.Alarm("CleanOrder error:", err)
		return err
	}
	order.CloseOrdId = orderId
	if err := order.watchCleanOrder(time.After(timeout)); err != nil {
		return err
	}
	return nil
}

func (order *Order) waitForOrderStart(timeout time.Duration, interval time.Duration) error {
	for {
		select {
		case <-time.After(timeout):
			log.Alarm("CleanOrder timeout, wait for order start")
			return &mo_errors.TimeoutError{}
		default:
		}
		if order.IsStart {
			break
		}
		select {
		case <-time.After(interval):
			continue
		}
	}
	return nil
}

func (order *Order) watchCleanOrder(wait <-chan time.Time) error {
	for {
		select {
		case <-wait:
			log.Alarm("CleanOrder timeout!!")
			return &mo_errors.TimeoutError{}
		default:
		}
		info, err := order.Trade.Market.GetOrderInfo(order.InstId, order.CloseOrdId)
		if err != nil {
			log.Println("GetOrderInfo error:", err)
			continue
		}
		if info.State == OKXClient.FILLED {
			order.Profit += info.Fee
			order.Profit += info.Pnl
			order.PriceOut = info.AvgPx
			order.IsFinished = true
			return nil
		}
	}
}

func (order *Order) Watching(interval time.Duration) {
	order.IsWatching = true
	for {
		t := time.NewTimer(interval)
		log.DebugStruct(order)
		if order.IsFinished {
			return
		}
		if !order.IsStart {
			select {
			case <-t.C:
				continue
			}
		}
		fee := order.Profit * 2
		v, err := order.Trade.Market.GetTickerValue(order.InstId)
		if err != nil {
			log.Println("GetTickerValue error:", err)
			continue
		}
		var profit float64
		if order.PosSide == OKXClient.LONG {
			priceDiff := v - order.PriceIn
			profit = priceDiff*float64(order.Size)*order.UnitSize + fee
		} else {
			priceDiff := order.PriceIn - v
			profit = priceDiff*float64(order.Size)*order.UnitSize + fee
		}
		if profit > 0 {
			order.IsProfitable = true
		} else {
			order.IsProfitable = false
		}
		log.Debugln("PriceIn:", order.PriceIn, "PriceNow", v, "Fee", fee, "Size", order.Size, "Unit", order.UnitSize)
		log.Debugln("Now profit:", profit)
		order.ApprProfit = profit
		select {
		case <-t.C:
			continue
		}
	}
}

func (order *Order) Protect(maxLoss float64, interval time.Duration) {
	order.IsProtect = true
	if !order.IsWatching {
		go order.Watching(interval)
	}
	for {
		t := time.NewTimer(interval)
		profit := order.ApprProfit
		if -1*profit > maxLoss {
			log.Println("Order", order.OpenOrdId, "[", order.InstId, "]", "LOSS PROTECTED!!!")
			if err := order.CleanOrder(OKXClient.MARKET, 0, time.Second*5); err != nil {
				if !order.IsFinished {
					log.Alarm("Clean FAILED, in LOSS PROTECTED!!!")
				}
				log.Debugln("Clean FAILED, But Finished...")
				return
			}
			log.Debugln("Clean SUCCESS")
			log.DebugStruct(order)
			return
		}
		select {
		case <-t.C:
			continue
		}
	}
}
