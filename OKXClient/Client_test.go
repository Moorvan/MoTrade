package OKXClient

import (
	"MoTrade/core"
	"MoTrade/global"
	"testing"
)

var (
	config *APIConfig
	client *OKX
)

func init() {
	core.Viper("../config_sim.yaml")
	config = &APIConfig{
		ApiKey:      global.GB_CONFIG.Api.ApiKey,
		SecretKey:   global.GB_CONFIG.Api.ApiSecretKey,
		Passphrase:  global.GB_CONFIG.Api.Passphrase,
		IsSimulated: global.GB_CONFIG.Api.SimulatedTrading,
	}
	client = New(config)
}

func TestRequestAllBalance(t *testing.T) {
	data, err := client.Account.GetAllBalance()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(data)
}

func TestRequestOneBalance(t *testing.T) {
	data, err := client.Account.GetOneBalance(USDT)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(data)
}

func TestRequestBalance(t *testing.T) {
	data, err := client.Account.GetBalance([]string{BTC, USDT})
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(data)
}
