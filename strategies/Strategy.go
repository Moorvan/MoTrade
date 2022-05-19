package strategies

import (
	"MoTrade/OKXClient"
	mlog "MoTrade/mo-log"
)

var (
	log = mlog.Log
)

type Strategy struct {
	Trade  *OKXClient.Trade
	Orders []Order
	InstId string
	TdMode string
}
