package cexioapi

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"reflect"
)


const (
	APIEndpoint = "https://cex.io/api/"
)


type CexioAPI struct {
	Username	string
	Key			string
	Secret		string
	Debug		bool
	Client		*http.Client
}

func NewCexioAPI(username string, key string, secret string, debug bool) (*CexioAPI){
	api := &CexioAPI{
		Username: 	username,
		Key:		key,
		Secret:		secret,
		Debug:		debug,
	}

	api.Client = &http.Client{}
	
	return api
}


func (api *CexioAPI) debugLog(ctx string, req *http.Request, res *http.Response) {
	if api.Debug {
		log.Println("### Debug ###")
		log.Printf("Context: %s", ctx)
		log.Printf("Request:\n%v\n", req)
		log.Printf("Response:\n%v", res)
		log.Println("#############")
	}
}

func (api *CexioAPI) APICall(endpoint string, method string, params string, data map[string][]string, private bool) (interface{}, error){
	u := APIEndpoint + endpoint

	if params != "" { u = u + "/" + params }
	
	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return nil, &customError{err, "request"}
	}
	
	if method == "POST" {
		v := url.Values(data)
		
		
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(v.Encode())))
		req.Body = ioutil.NopCloser(strings.NewReader(v.Encode()))
	}
	
	res, err := api.Client.Do(req)
	if err != nil {
		return nil, &customError{err, "response"}
	}

	if res.StatusCode != 200 {
		e := fmt.Errorf("HTTP Error Code: %d", res.StatusCode)
		return nil, &customError{e, "response"}
	}

	defer res.Body.Close()

	api.debugLog(endpoint, req, res)

	var t interface{}
	err = json.NewDecoder(res.Body).Decode(&t)
	if err != nil {
		return nil, &customError{err, "(JSON) Decode response"}
	}

	// API does not return a standard format
	typeOf := reflect.ValueOf(t)
	if typeOf.Kind() == reflect.Map {
		data := t.(map[string]interface{})
		err, exist := data["error"]
		if exist {
			e := fmt.Errorf("%s", err)
			return nil, &customError{e, "API error"}
		}
	}

	return t, nil
}


//Public functions

func (api *CexioAPI) CurrencyLimits() (map[string]interface{}, error){
	res, err := api.APICall("currency_limits", "GET", "", nil, false)
	return res.(map[string]interface{}), err
}

func (api *CexioAPI) Ticker(base string, currency string) (map[string]interface{}, error){
	params := strings.ToUpper(base + "/" + currency)
	res, err := api.APICall("ticker", "GET", params, nil, false)
	return res.(map[string]interface{}), err
}

func (api *CexioAPI) Tickers(base string, currency string) (map[string]interface{}, error){
	params := strings.ToUpper(base + "/" + currency)
	res, err := api.APICall("tickers", "GET", params, nil, false)
	return res.(map[string]interface{}), err
}

func (api *CexioAPI) LastPrice(base string, currency string) (interface{}, error){
	params := strings.ToUpper(base + "/" + currency)
	res, err := api.APICall("last_price", "GET", params, nil, false)
	return res.(map[string]interface{}), err
}

func (api *CexioAPI) LastPrices(base string, currency string, currency2 string) (interface{}, error){
	params := strings.ToUpper(base + "/" + currency + "/" + currency2)
	res, err := api.APICall("last_prices", "GET", params, nil, false)
	return res.(map[string]interface{}), err
}

func (api *CexioAPI) Converter(base string, currency string, amount string) (interface{}, error){
	params := strings.ToUpper(base + "/" + currency)
	data := map[string][]string{ "amnt": []string{amount} }
	res, err := api.APICall("convert", "POST", params, data, false)
	return res.(map[string]interface{}), err
}

func (api *CexioAPI) Chart(base string, currency string, lastHours string, maxRespArrSize string) ([]interface{}, error){
	params := strings.ToUpper(base + "/" + currency)
	data := map[string][]string {
		"lastHours": []string{lastHours},
		"maxRespArrSize": []string{maxRespArrSize},
	}
	res, err := api.APICall("price_stats", "POST", params, data, false)
	return res.([]interface{}), err
}

func (api *CexioAPI) OhlcvChart(base string, currency string, date string) (interface{}, error){
	params := date + "/" + strings.ToUpper(base + "/" + currency)
	res, err := api.APICall("ohlcv/hd", "GET", params, nil, false)
	return res.(map[string]interface{}), err
}

func (api *CexioAPI) Orderbook(base string, currency string) (interface{}, error){
	params := strings.ToUpper(base + "/" + currency)
	res, err := api.APICall("order_book", "GET", params, nil, false)
	return res.(map[string]interface{}), err
}

func (api *CexioAPI) TradeHistory(base string, currency string) ([]interface{}, error){
	params := strings.ToUpper(base + "/" + currency)
	res, err := api.APICall("trade_history", "GET", params, nil, false)
	return res.([]interface{}), err
}

