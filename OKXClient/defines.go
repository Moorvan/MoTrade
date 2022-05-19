package OKXClient

// Currency
const (
	DOGE = "DOGE"
	BTC  = "BTC"
	USDT = "USDT"
)

// InstType
const (
	SWAP    = "SWAP"
	SPOT    = "SPOT"
	MARGIN  = "MARGIN"
	FUTURES = "FUTURES"
)

// InstID
const (
	DOGE_USDT_SWAP = "DOGE-USDT-SWAP"
	BTC_USDT_SWAP  = "BTC-USDT-SWAP"
	ETH_USDT_SWAP  = "ETH-USDT-SWAP"
)

// Duration
const (
	MINUTE_1  = "1m"
	MINUTE_3  = "3m"
	MINUTE_5  = "5m"
	MINUTE_15 = "15m"
	MINUTE_30 = "30m"
	HOUR_1    = "1h"
	HOUR_2    = "2h"
	HOUR_4    = "4h"
)

// Trade Mode
const (
	ISOLATED = "isolated"
	CROSS    = "cross"
	CASH     = "cash"
)

// Side
const (
	BUY  = "buy"
	SELL = "sell"
)

// position side
const (
	LONG  = "long"
	SHORT = "short"
	NET   = "net"
)

// Order Type
const (
	LIMIT  = "limit"
	MARKET = "market"
)

// Order State
const (
	CANCELED         = "canceled"
	LIVE             = "live"
	PARTIALLY_FILLED = "partially_filled"
	FILLED           = "filled"
)
