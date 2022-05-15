package tests

import (
	"MoTrade/OKXClient"
	"MoTrade/core"
	mlog "MoTrade/mo-log"
	"testing"
)

var (
	config *OKXClient.APIConfig
	client *OKXClient.OKX
)

var (
	log = mlog.Log
)

func init() {
	core.Viper("../config_sim.yaml")
	client = core.NewOKX()
}

func TestTickerApi(t *testing.T) {
	api := "/api/v5/market/ticker"
	response := &struct {
		Data any
	}{}
	params := OKXClient.ParamsBuilder().Set("instId", OKXClient.DOGE_USDT_SWAP)
	if err := client.DoGet(api, params, response); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(response)
}

func TestGetTickerValue(t *testing.T) {
	data, err := client.Market.GetTickerValue(OKXClient.DOGE_USDT_SWAP)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(data)
}

func TestRequestMA(t *testing.T) {
	data, err := client.Market.GetMA(OKXClient.DOGE_USDT_SWAP, OKXClient.MINUTE_1, 20, 1)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(data)
}

func TestRequestAllBalance(t *testing.T) {
	data, err := client.Account.GetAllBalance()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(data)
}

func TestRequestOneBalance(t *testing.T) {
	data, err := client.Account.GetOneBalance(OKXClient.USDT)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(data)
}

func TestRequestBalance(t *testing.T) {
	data, err := client.Account.GetBalance([]string{OKXClient.BTC, OKXClient.USDT})
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(data)
}
