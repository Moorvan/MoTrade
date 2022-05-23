package OKXClient

import (
	mo_errors "MoTrade/mo-errors"
	"errors"
	"fmt"
	"strings"
)

type AccountAPI interface {
	GetBalance(ccy []string) (map[string]float64, error)
	GetOneBalance(ccy string) (float64, error)
	GetAllBalance() (map[string]float64, error)
}

type OKXAccountAPI struct {
	*OKX
}

func (account *OKXAccountAPI) GetBalance(ccy []string) (map[string]float64, error) {
	api := "/api/v5/account/balance"
	currencys := strings.Join(ccy, ",")
	params := ParamsBuilder().Set("ccy", currencys)

	response := &struct {
		Data []struct {
			Details []struct {
				Ccy string
				Eq  float64 `json:"eq,string"`
			}
		}
	}{}
	err := account.DoGet(api, params, response)
	if err != nil {
		log.Errorln("GetBalance of " + api + " Failed: " + err.Error())
		return nil, errors.New("GetBalance of " + api + " Failed: " + err.Error())
	}
	log.Printf("%+v", response)

	balance := map[string]float64{}
	for _, v := range response.Data[0].Details {
		balance[v.Ccy] = v.Eq
	}
	return balance, err
}

func (account *OKXAccountAPI) GetOneBalance(ccy string) (float64, error) {
	res, err := account.GetBalance([]string{ccy})
	if err != nil {
		log.Errorln("GetOneBalance of", ccy, "Fail:", err.Error())
		return 0, err
	}
	if r, ok := res[ccy]; ok {
		return r, nil
	} else {
		return 0, mo_errors.NoResultError
	}
}

type Response struct {
	Code int `json:"code,string"`
	Msg  string
}

func (account *OKXAccountAPI) GetAllBalance() (map[string]float64, error) {
	api := "/api/v5/account/balance"
	response := struct {
		Response
		Data []struct {
			Details []struct {
				Ccy string
				Eq  float64 `json:"eq,string"`
			}
		}
	}{}
	err := account.DoGet(api, nil, &response)
	if err != nil {
		log.Errorln("GetBalance of " + api + " Failed: " + err.Error())
		return nil, errors.New("GetBalance of " + api + " Failed: " + err.Error())
	}
	fmt.Println(response)

	log.Printf("%+v", response)
	balance := map[string]float64{}
	for _, v := range response.Data[0].Details {
		balance[v.Ccy] = v.Eq
	}
	return balance, err
}
