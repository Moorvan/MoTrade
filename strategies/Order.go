package strategies

import (
	"MoTrade/OKXClient"
	mo_errors "MoTrade/mo-errors"
	"time"
)

type Order struct {
	Trade      *OKXClient.Trade
	InstId     string
	TdMode     string
	OrdType    string
	PosSide    string
	OpenOrdId  string
	CloseOrdId string
	PriceIn    float64
	PriceOut   float64
	Size       int
	IsStart    bool
	IsFinished bool
	Profit     float64
}

func NewOrder(trade *OKXClient.Trade, instId, tdMode, posSide, ordType string, size int, px float64, timeout time.Duration) (*Order, error) {
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
	order := &Order{
		Trade:     trade,
		InstId:    instId,
		TdMode:    tdMode,
		OrdType:   ordType,
		PosSide:   posSide,
		OpenOrdId: ordId,
		Size:      size,
		Profit:    0,
	}
	if err := order.watchOpenOrder(time.After(timeout)); err != nil {
		return order, err
	}
	return order, nil
}

func (order *Order) watchOpenOrder(wait <-chan time.Time) error {
	for {
		select {
		case <-wait:
			if err := order.CancelOrder(); err != nil {
				log.Alarm("CancelOrder error:", err)
				return err
			}
			info, err := order.Trade.Market.GetOrderInfo(order.InstId, order.OpenOrdId)
			if err != nil {
				log.Fatalln("GetOrderInfo error:", err)
				return err
			}
			order.Size = info.Size
			order.PriceIn = info.AvgPx
			order.Profit += info.Fee
			if order.Size == 0 {
				order.IsFinished = true
			}
			order.IsStart = true
			return &mo_errors.TimeoutError{}
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
			order.IsStart = true
			return nil
		}
	}
}

func (order *Order) CancelOrder() error {
	if err := order.Trade.Market.CancelOrder(order.InstId, order.OpenOrdId); err != nil {
		log.Println("CancelOrder error:", err)
		return err
	}
	return nil
}

func (order *Order) CleanOrder(ordType string, px float64, timeout time.Duration) error {
	var side string
	if order.PosSide == OKXClient.LONG {
		side = OKXClient.SELL
	} else {
		side = OKXClient.BUY
	}
	orderId, err := order.Trade.Market.PlaceOrder(order.InstId, order.TdMode, side, order.PosSide, ordType, order.Size, px)
	if err != nil {
		log.Errorln("CleanOrder error:", err)
		return err
	}
	order.CloseOrdId = orderId
	if err := order.watchCleanOrder(time.After(timeout)); err != nil {
		return err
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
