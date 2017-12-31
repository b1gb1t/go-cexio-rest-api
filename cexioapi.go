package cexioapi

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"encoding/json"
	"io/ioutil"
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

func checkError(ctx string, err error) {
	if err != nil {
		log.Fatalf("Context: %s; Error: %v", ctx, err)
	}
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

func (api *CexioAPI) APICall(endpoint string, method string, params string, data map[string][]string, private bool) (interface{}){
	u := APIEndpoint + endpoint

	if params != "" { u = u + "/" + params }
	
	req, err := http.NewRequest(method, u, nil)
	checkError("request", err)
	
	if method == "POST" {
		v := url.Values(data)
		
		
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(v.Encode())))
		req.Body = ioutil.NopCloser(strings.NewReader(v.Encode()))
	}
	
	res, err := api.Client.Do(req)
	checkError("response", err)
	defer res.Body.Close()

	api.debugLog(endpoint, req, res)

	var t interface{}
	err = json.NewDecoder(res.Body).Decode(&t)
	checkError("json - parse response", err)

	return t
}


//Public functions

func (api *CexioAPI) CurrencyLimits() (interface{}){
	return api.APICall("currency_limits", "GET", "", nil, false)
}

func (api *CexioAPI) Ticker(base string, currency string) (interface{}){
	params := strings.ToUpper(base + "/" + currency)
	return api.APICall("ticker", "GET", params, nil, false)
}

func (api *CexioAPI) Tickers(base string, currency string) (interface{}){
	params := strings.ToUpper(base + "/" + currency)
	return api.APICall("tickers", "GET", params, nil, false)
}

func (api *CexioAPI) LastPrice(base string, currency string) (interface{}){
	params := strings.ToUpper(base + "/" + currency)
	return api.APICall("last_price", "GET", params, nil, false)
}

func (api *CexioAPI) LastPrices(base string, currency string, currency2 string) (interface{}){
	params := strings.ToUpper(base + "/" + currency + "/" + currency2)
	return api.APICall("last_prices", "GET", params, nil, false)
}

func (api *CexioAPI) Converter(base string, currency string, amount string) (interface{}){
	params := strings.ToUpper(base + "/" + currency)
	data := map[string][]string{ "amnt": []string{amount} }
	return api.APICall("convert", "POST", params, data, false)
}

func (api *CexioAPI) Chart(base string, currency string, lastHours string, maxRespArrSize string) (interface{}){
	params := strings.ToUpper(base + "/" + currency)
	data := map[string][]string {
		"lastHours": []string{lastHours},
		"maxRespArrSize": []string{maxRespArrSize},
	}
	return api.APICall("price_stats", "POST", params, data, false)
}

func (api *CexioAPI) OhlcvChart(base string, currency string, date string) (interface{}){
	params := date + "/" + strings.ToUpper(base + "/" + currency)
	return api.APICall("ohlcv/hd", "GET", params, nil, false)
}

func (api *CexioAPI) Orderbook(base string, currency string) (interface{}){
	params := strings.ToUpper(base + "/" + currency)
	return api.APICall("order_book", "GET", params, nil, false)
}

func (api *CexioAPI) TradeHistory(base string, currency string) (interface{}){
	params := strings.ToUpper(base + "/" + currency)
	return api.APICall("trade_history", "GET", params, nil, false)

}

