package strategies

import (
	"MoTrade/OKXClient"
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

func NewOrder(trade *OKXClient.Trade, instId, tdMode, posSide, ordType string, size int, px float64) (*Order, error) {
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
	}
	// TODO: Open success?
	return order, nil
}

func (order *Order) CancelOrder() error {
	if err := order.Trade.Market.CancelOrder(order.InstId, order.OpenOrdId); err != nil {
		log.Println("CancelOrder error:", err)
		return err
	}
	// TODO: Cancel Success?
	return nil
}

func (order *Order) CleanOrder(ordType string, px float64) error {
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
	// TODO: Clean Success? Profit calculation.
	return nil
}
