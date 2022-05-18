package tests

import (
	"MoTrade/OKXClient"
	"testing"
	"time"
)

func TestGetOrdersHistoryApi(t *testing.T) {
	api := "/api/v5/account/bills"
	params := OKXClient.ParamsBuilder().Set("instType", OKXClient.SWAP)
	response := &struct {
		Data []struct {
			OrdId  string  `json:"ordId"`
			Pnl    float64 `json:"pnl,string"`
			Sz     int     `json:"sz,string"`
			Type   string  `json:"type"`
			BalChg float64 `json:"balChg,string"`
		}
	}{}

	if err := client.DoGet(api, params, response); err != nil {
		t.Error(err)
	}
	log.Println(response)
}

func TestPlaceOneOrderAndCancelOrder(t *testing.T) {
	orderId, err := client.Market.PlaceOrder(OKXClient.ETH_USDT_SWAP, OKXClient.CROSS, OKXClient.BUY, OKXClient.LONG, OKXClient.MARKET, 10, 0)
	if err != nil {
		t.Error(err)
	}
	log.Println(orderId)
	if err := client.Market.CancelOrder(OKXClient.ETH_USDT_SWAP, orderId); err != nil {
		log.Println("cancel order error:", err)
	}
	info, err := client.Market.GetOrderInfo(OKXClient.ETH_USDT_SWAP, orderId)
	if err != nil {
		t.Error(err)
	}
	log.Println(info)
}

func TestPlaceOneOrderAndClosePosition(t *testing.T) {
	orderId, err := client.Market.PlaceOrder(OKXClient.ETH_USDT_SWAP, OKXClient.CROSS, OKXClient.BUY, OKXClient.LONG, OKXClient.MARKET, 10, 0)
	if err != nil {
		t.Error(err)
	}
	log.Println(orderId)
	time.Sleep(1 * time.Second)
	if err := client.Market.ClosePosition(OKXClient.ETH_USDT_SWAP, OKXClient.LONG, OKXClient.CROSS); err != nil {
		log.Println("close position error:", err)
		t.Error(err)
	}
}

func TestPlaceOneOrderAndClossPosition2(t *testing.T) {
	// Buy Long -> Sell Long
	// Sell Short -> Buy Short
	orderId, err := client.Market.PlaceOrder(OKXClient.ETH_USDT_SWAP, OKXClient.CROSS, OKXClient.SELL, OKXClient.SHORT, OKXClient.MARKET, 10, 0)
	if err != nil {
		t.Error(err)
	}
	log.Println(orderId)
	info, err := client.Market.GetOrderInfo(OKXClient.ETH_USDT_SWAP, orderId)
	if err != nil {
		t.Error(err)
	}
	log.Println(info)

	log.Println("Wait 10 seconds")
	time.Sleep(10 * time.Second)

	orderId, err = client.Market.PlaceOrder(OKXClient.ETH_USDT_SWAP, OKXClient.CROSS, OKXClient.BUY, OKXClient.SHORT, OKXClient.MARKET, 10, 0)
	if err != nil {
		t.Error(err)
	}
	log.Println(orderId)
	info, err = client.Market.GetOrderInfo(OKXClient.ETH_USDT_SWAP, orderId)
	if err != nil {
		t.Error(err)
	}
	log.Println(info)
}

func TestRequestGetOrder(t *testing.T) {
	data, err := client.Market.GetOrderInfo(OKXClient.ETH_USDT_SWAP, "447209241030561793")
	if err != nil {
		t.Error(err)
	}
	log.Println(data)
}

func TestGetOrderApi(t *testing.T) {
	api := "/api/v5/trade/order"

	params := OKXClient.ParamsBuilder().Set("instId", OKXClient.ETH_USDT_SWAP).Set("ordId", "447162784663605248")
	response := &struct {
		Data []struct {
			Pnl   string  `json:"pnl"`
			State string  `json:"state"`
			AvgPx float64 `json:"avgPx,string"`
			Sz    int     `json:"sz,string"`
			Fee   float64 `json:"fee,string"`
		}
	}{}
	if err := client.DoGet(api, params, response); err != nil {
		t.Error(err)
	}
	log.Println(response)

	data, err := client.Market.GetTickerValue(OKXClient.ETH_USDT_SWAP)
	if err != nil {
		t.Error(err)
	}
	log.Println(data)
}

func TestOrderApi(t *testing.T) {
	api := "/api/v5/trade/order"
	request := &struct {
		InstId  string
		TdMode  string
		Side    string
		PosSide string
		OrdType string
		Sz      string
		Px      string
	}{
		InstId:  OKXClient.DOGE_USDT_SWAP,
		TdMode:  OKXClient.CROSS,
		Side:    OKXClient.BUY,
		PosSide: OKXClient.LONG,
		OrdType: OKXClient.MARKET,
		Sz:      "1",
		Px:      "0.08",
	}

	response := &struct {
		Data any
	}{}
	if err := client.DoPost(api, request, response); err != nil {
		t.Error(err)
	}
	//log.Println(response.Data)
	log.Println(response)
}

func TestRequestPlaceOneOrder(t *testing.T) {
	data, err := client.Market.PlaceOrder(OKXClient.ETH_USDT_SWAP, OKXClient.CROSS, OKXClient.BUY, OKXClient.LONG, OKXClient.MARKET, 10, 0)
	if err != nil {
		t.Error(err)
	}
	log.Println(data)
}

func TestCancelOrderApi(t *testing.T) {
	api := "/api/v5/trade/cancel-order"
	request := &struct {
		InstId string
		OrdId  string
	}{
		InstId: OKXClient.ETH_USDT_SWAP,
		OrdId:  "446091617966100480",
	}
	response := &struct {
		Data any
	}{}
	if err := client.DoPost(api, request, response); err != nil {
		t.Error(err)
	}
	log.Println(response.Data)
}

func TestRequestCancelOrder(t *testing.T) {
	err := client.Market.CancelOrder(OKXClient.DOGE_USDT_SWAP, "446098262561525760")
	if err != nil {
		t.Error(err)
	}
}

func TestClosePositionApi(t *testing.T) {
	api := "/api/v5/trade/close-position"
	request := &struct {
		InstId  string
		PosSide string
		MgnMode string
	}{
		InstId:  OKXClient.ETH_USDT_SWAP,
		PosSide: OKXClient.LONG,
		MgnMode: OKXClient.CROSS,
	}
	response := &struct {
		Data any
	}{}
	if err := client.DoPost(api, request, response); err != nil {
		t.Error(err)
	}
	log.Println(response.Data)
}

func TestRequestCancelPosition(t *testing.T) {
	err := client.Market.ClosePosition(OKXClient.DOGE_USDT_SWAP, OKXClient.LONG, OKXClient.CROSS)
	if err != nil {
		t.Error(err)
	}
}
