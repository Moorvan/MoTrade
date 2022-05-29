package main

import (
	"MoTrade/OKXClient"
	"MoTrade/core"
	"MoTrade/global"
	mlog "MoTrade/mo-log"
	"MoTrade/strategies/swap/MABased"
	"time"
)

var (
	log        = mlog.Log
	configPath = "config_sim.yaml"
)

func main() {
	global.GB_VP = core.Viper(configPath)
	global.GB_CLIENT = core.NewOKX()
	strategy := MABased.NewMABasedStrategy(&global.GB_CLIENT.Trade, 3, OKXClient.SWAP, OKXClient.ETH_USDT_SWAP, OKXClient.CROSS, 1, OKXClient.MINUTE_1, 5)
	strategy.Run(time.Second)

}
