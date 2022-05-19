package tests

import (
	"MoTrade/OKXClient"
	"MoTrade/strategies"
	"testing"
	"time"
)

var (
	order *strategies.Order
)

func PlaceOneOrder() {
	var err error
	order, err = strategies.NewOrder(&client.Trade, OKXClient.ETH_USDT_SWAP, OKXClient.CROSS, OKXClient.LONG, OKXClient.MARKET, 1, 0, 5*time.Second)
	if err != nil {
		panic(err)
	}
}

func TestCancelOrder(t *testing.T) {
	PlaceOneOrder()
	if err := order.CancelOrder(); err != nil {
		log.Println(err)
	}
}

func TestCleanOrder(t *testing.T) {
	PlaceOneOrder()
	time.Sleep(5 * time.Second)
	if err := order.CleanOrder(OKXClient.MARKET, 0, 5*time.Second); err != nil {
		log.Errorln(err)
	}

	log.Printf("%+v", order)

}
