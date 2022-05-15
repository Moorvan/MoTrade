package test

import (
	"MoTrade/OKXClient"
	"testing"
)

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
	data, err := client.Market.PlaceOrder(OKXClient.DOGE_USDT_SWAP, OKXClient.CROSS, OKXClient.BUY, OKXClient.LONG, OKXClient.MARKET, 1, 0.08)
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
		InstId: OKXClient.DOGE_USDT_SWAP,
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
		InstId:  OKXClient.DOGE_USDT_SWAP,
		PosSide: OKXClient.SHORT,
		MgnMode: OKXClient.ISOLATED,
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
