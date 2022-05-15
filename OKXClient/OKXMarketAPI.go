package OKXClient

import "strconv"

type MarketAPI interface {
	GetMA(instId, bar string, limit int) (float64, error)
}

type OKXMarketAPI struct {
	*OKX
}

func (market OKXMarketAPI) GetMA(instId, bar string, limit int) (float64, error) {
	api := "/api/v5/market/candles"
	params := ParamsBuilder().Set("instId", instId).Set("bar", bar).Set("limit", strconv.Itoa(limit))

	response := &struct {
		Data [][]string
	}{}
	if err := market.DoGet(api, params, response); err != nil {
		log.Errorln("GetMA error:", err.Error())
		return 0, err
	}

	log.Println(response)
	log.Println(len(response.Data))
	l := len(response.Data)
	sum := 0.0
	for _, v := range response.Data {
		c, err := strconv.ParseFloat(v[4], 64)
		if err != nil {
			log.Errorln("GetMA data error:", err.Error())
			return 0, err
		}
		sum += c
	}
	return sum / float64(l), nil
}
