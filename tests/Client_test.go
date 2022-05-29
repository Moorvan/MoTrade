package tests

import (
	"MoTrade/OKXClient"
	"MoTrade/core"
	mlog "MoTrade/mo-log"
	"testing"
	"time"
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

func TestDuration(t *testing.T) {
	var err error
	_, err = time.ParseDuration(OKXClient.MINUTE_1)
	if err != nil {
		t.Error(err)
	}
	_, err = time.ParseDuration(OKXClient.MINUTE_3)
	if err != nil {
		t.Error(err)
	}
	_, err = time.ParseDuration(OKXClient.MINUTE_5)
	if err != nil {
		t.Error(err)
	}
	_, err = time.ParseDuration(OKXClient.MINUTE_15)
	if err != nil {
		t.Error(err)
	}
	_, err = time.ParseDuration(OKXClient.MINUTE_30)
	if err != nil {
		t.Error(err)
	}
	_, err = time.ParseDuration(OKXClient.HOUR_1)
	if err != nil {
		t.Error(err)
	}
	_, err = time.ParseDuration(OKXClient.HOUR_2)
	if err != nil {
		t.Error(err)
	}
	_, err = time.ParseDuration(OKXClient.HOUR_4)
	if err != nil {
		t.Error(err)
	}
}

func TestTickerApi(t *testing.T) {
	api := "/api/v5/market/ticker"
	response := &struct {
		Data any
	}{}
	params := OKXClient.ParamsBuilder().Set("instId", OKXClient.ETH_USDT_SWAP)
	if err := client.DoGet(api, params, response); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(response)
}

func TestTickerMsgApi(t *testing.T) {
	api := "/api/v5/public/instruments"
	params := OKXClient.ParamsBuilder().Set("instType", OKXClient.SWAP).Set("instId", OKXClient.ETH_USDT_SWAP)
	response := &struct {
		Data []struct {
			LotSz string `json:"lotSz"`
			MinSz string `json:"minSz"`
			CtVal string `json:"ctVal"`
		}
	}{}

	if err := client.DoGet(api, params, response); err != nil {
		log.Fatalln(err.Error())
	}
	log.PrintStruct(response)
}

func TestRequestForTickerUnitSize(t *testing.T) {
	ret, err := client.Market.GetTickerUnitSize(OKXClient.SWAP, OKXClient.ETH_USDT_SWAP)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(ret)
}

func TestGetTickerValue(t *testing.T) {
	data, err := client.Market.GetTickerValue(OKXClient.ETH_USDT_SWAP)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(data)
}

func TestRequestMA(t *testing.T) {
	data, err := client.Market.GetMA(OKXClient.ETH_USDT_SWAP, OKXClient.MINUTE_1, 5, 1)
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
