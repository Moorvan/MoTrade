package main

import (
	"MoTrade/core"
	"MoTrade/global"
	mlog "MoTrade/mo-log"
)

var log = mlog.Log

func main() {
	global.GB_VP = core.Viper("config_sim.yaml")
}
