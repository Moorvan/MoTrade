package OKXClient

import (
	mlog "MoTrade/mo-log"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"time"
)

var (
	baseUrl = "https://www.okx.com"
	log     = mlog.Log
)

type Trade struct {
	Account AccountAPI
	Market  MarketAPI
}

type OKX struct {
	config *APIConfig
	client *resty.Client
	Trade
}

func New(conf *APIConfig) *OKX {
	client := resty.New()
	client.SetHeaders(map[string]string{
		"Content-Type":         "application/json",
		"OK-ACCESS-KEY":        conf.ApiKey,
		"OK-ACCESS-PASSPHRASE": conf.Passphrase,
	})
	log.Println("IsSimulated:", conf.IsSimulated)
	if conf.IsSimulated {
		client.SetHeader("x-simulated-trading", "1")
	}

	okx := &OKX{
		config: conf,
		client: client,
	}
	okx.Account = &OKXAccountAPI{okx}
	okx.Market = &OKXMarketAPI{okx}

	return okx
}

func (okx *OKX) DoGet(api string, params Params, response any) error {
	req := okx.client.R()

	apiWithParams := api
	if params != nil {
		req.SetQueryParams(params)
		apiWithParams = api + "?" + req.QueryParam.Encode()
	}
	okx.updateHeadersWithNowTime(req, "GET", apiWithParams)

	resp, err := req.SetResult(response).Get(baseUrl + api)

	if err := checkResponse(resp); err != nil {
		return err
	}

	//log.Println(resp)
	if err != nil {
		log.Errorln("Get Request", api, "FAIL", err.Error())
		return err
	}

	return nil
}

func (okx *OKX) DoPost(api string, body, response any) error {
	req := okx.client.R()
	s, err := json.Marshal(body)
	if err != nil {
		log.Errorln("Marshal Request", api, "FAIL", err.Error())
		return err
	}
	okx.updateHeadersWithNowTime(req, "POST", api+string(s))
	//log.Println("Post Headers:", req.Header)
	//log.Println("Post Headers:", okx.client.Header)

	resp, err := req.SetBody(body).SetResult(response).Post(baseUrl + api)
	log.Println(resp)

	if err := checkResponse(resp); err != nil {
		return err
	}

	if err != nil {
		log.Errorln("Post Request", api, "FAIL:", err.Error())
		return err
	}

	return nil
}

func checkResponse(resp *resty.Response) error {
	status := &struct {
		Code int `json:"code,string"`
		Msg  string
	}{}
	if err := json.Unmarshal([]byte(resp.String()), status); err != nil {
		log.Fatalln("Can't connect")
	}
	if status.Code != 0 {
		return errors.New("Response with Fail: " + status.Msg)
	}
	return nil
}

func (okx *OKX) updateHeadersWithNowTime(req *resty.Request, method, apiWithParams string) {
	timeStamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	s := timeStamp + method + apiWithParams
	h := hmac.New(sha256.New, []byte(okx.config.SecretKey))
	h.Write([]byte(s))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	req.SetHeaders(map[string]string{
		"OK-ACCESS-SIGN":      sign,
		"OK-ACCESS-TIMESTAMP": timeStamp,
	})
}
