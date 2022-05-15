package global

import (
	"MoTrade/OKXClient"
	"MoTrade/config"
	"github.com/spf13/viper"
)

var (
	GB_VP     *viper.Viper
	GB_CONFIG config.Config
	GB_CLIENT *OKXClient.OKX
)
