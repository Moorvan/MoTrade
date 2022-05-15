package tests

import (
	"MoTrade/OKXClient"
	"MoTrade/core"
	mlog "MoTrade/mo-log"
	"strconv"
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

func TestApi(t *testing.T) {
	response := &struct {
		Data [][]string
	}{}
	ts := time.Now().Add(-5 * time.Minute).UnixMilli()
	params := OKXClient.ParamsBuilder().Set("instId", OKXClient.DOGE_USDT_SWAP).Set("limit", "20").Set("after", strconv.FormatInt(ts, 10))
	if err := client.DoGet("/api/v5/market/candles", params, response); err != nil {
		log.Fatalln(err.Error())
	}
	//log.Println(response.Data)
	for _, v := range response.Data {
		t, _ := strconv.Atoi(v[0])
		log.Println(time.UnixMilli(int64(t)))
		v, _ := strconv.ParseFloat(v[4], 64)
		log.Println(v)
	}
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
