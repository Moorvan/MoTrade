package OKXClient

import (
	"errors"
	"strconv"
	"time"
)

type MarketAPI interface {
	GetMA(instId, bar string, limit int, before time.Duration) (float64, error)
	PlaceOrder(instId, tdMode, side, posSide, ordType string, size int, px float64) (ordId string, err error)
	CancelOrder(instId, ordId string) error
	ClosePosition(instId, posSide, tradeMode string) error
}

type OKXMarketAPI struct {
	*OKX
}

func (market *OKXMarketAPI) GetMA(instId, bar string, limit int, before time.Duration) (float64, error) {
	api := "/api/v5/market/candles"
	ts := time.Now().Add(-1 * before).UnixMilli()
	params := ParamsBuilder().Set("instId", instId).Set("bar", bar).Set("limit", strconv.Itoa(limit)).Set("after", strconv.FormatInt(ts, 10))

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

func (market *OKXMarketAPI) PlaceOrder(instId, tdMode, side, posSide, ordType string, size int, px float64) (ordId string, err error) {
	api := "/api/v5/trade/order"
	request := &struct {
		InstId  string
		TdMode  string
		Side    string
		PosSide string
		OrdType string
		Sz      string
		Px      string
	}{
		InstId:  instId,
		TdMode:  tdMode,
		Side:    side,
		PosSide: posSide,
		OrdType: ordType,
		Sz:      strconv.Itoa(size),
		Px:      strconv.FormatFloat(px, 'f', -1, 64),
	}

	response := &struct {
		Data []struct {
			SCode string
			SMsg  string
			OrdId string
		}
	}{}

	if err := market.DoPost(api, request, response); err != nil {
		log.Errorln("PlaceOrder error:", err.Error())
		return "", err
	}
	if len(response.Data) == 0 {
		log.Errorln("PlaceOrder error:", errors.New("response data is empty"))
		return "", errors.New("PlaceOrder response data is empty")
	}
	data := response.Data[0]
	if data.SCode != "0" {
		log.Errorln("PlaceOrder error:", errors.New(data.SMsg))
		return "", errors.New(data.SMsg)
	}
	return data.OrdId, nil
}

func (market *OKXMarketAPI) CancelOrder(instId, ordId string) error {
	api := "/api/v5/trade/cancel-order"
	request := &struct {
		InstId string
		OrdId  string
	}{
		InstId: instId,
		OrdId:  ordId,
	}
	response := &struct {
		Data []struct {
			SCode string
			SMsg  string
		}
	}{}
	if err := market.DoPost(api, request, response); err != nil {
		log.Errorln("CancelOrder error:", err.Error())
		return err
	}
	if len(response.Data) == 0 {
		log.Errorln("CancelOrder error:", errors.New("response data is empty"))
		return errors.New("CancelOrder response data is empty")
	}
	data := response.Data[0]
	if data.SCode != "0" {
		log.Errorln("CancelOrder error:", errors.New(data.SMsg))
		return errors.New(data.SMsg)
	}
	return nil
}

func (market *OKXMarketAPI) ClosePosition(instId, posSide, tradeMode string) error {
	api := "/api/v5/trade/close-position"
	request := &struct {
		InstId  string
		PosSide string
		MgnMode string
	}{
		InstId:  instId,
		PosSide: posSide,
		MgnMode: tradeMode,
	}
	response := &struct {
		Data any
	}{}

	if err := market.DoPost(api, request, response); err != nil {
		log.Errorln("ClosePosition error:", err.Error())
		return err
	}
	return nil
}