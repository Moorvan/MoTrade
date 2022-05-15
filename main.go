package main

import (
	"MoTrade/core"
	"MoTrade/global"
	mlog "MoTrade/mo-log"
)

var (
	log        = mlog.Log
	configPath = "config_sim.yaml"
)

func main() {
	global.GB_VP = core.Viper(configPath)
	global.GB_Client = core.NewOKX()

}
