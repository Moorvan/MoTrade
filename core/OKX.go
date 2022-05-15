package core

import (
	"MoTrade/OKXClient"
	"MoTrade/global"
)

func NewOKX() *OKXClient.OKX {
	gbConfig := global.GB_CONFIG.Api
	config := &OKXClient.APIConfig{
		ApiKey:      gbConfig.ApiKey,
		SecretKey:   gbConfig.ApiSecretKey,
		Passphrase:  gbConfig.Passphrase,
		IsSimulated: gbConfig.SimulatedTrading,
	}
	return OKXClient.New(config)
}
