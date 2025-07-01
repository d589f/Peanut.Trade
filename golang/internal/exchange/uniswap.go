package exchange

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang/internal/marketdata"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type UniswapProvider struct{}

// ВСТАВИТЬ СВОЙ API
var uniswapV3GraphUrl = "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3"

var uniswapPools = map[string]string{
	//TODO
}

func (u *UniswapProvider) Name() string { return "Uniswap" }

func (u *UniswapProvider) GetLatestPrice(pair string) (float64, error) {
	poolId, ok := uniswapPools[strings.ToUpper(pair)]
	if !ok {
		return 0, fmt.Errorf("unsupported pair for Uniswap: %s", pair)
	}
	query := fmt.Sprintf(`{"query": "{ pool(id: \"%s\") { token0 { symbol } token1 { symbol } token0Price token1Price } }"}`, poolId)
	resp, err := http.Post(uniswapV3GraphUrl, "application/json", bytes.NewBuffer([]byte(query)))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("uniswap: bad status %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var gql struct {
		Data struct {
			Pool struct {
				Token0      struct{ Symbol string }
				Token1      struct{ Symbol string }
				Token0Price string
				Token1Price string
			}
		}
	}
	if err := json.Unmarshal(body, &gql); err != nil {
		return 0, err
	}
	// Для пары ETH/USDC: цена ETH в USDC = token1Price
	price, err := parseFloat(gql.Data.Pool.Token1Price)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func (u *UniswapProvider) Poll(pair string, md marketdata.MarketDataServiceUpdater) {
	price, err := u.GetLatestPrice(pair)
	if err != nil {
		md.UpdateError(u.Name(), pair, err.Error())
		return
	}
	md.UpdatePrice(u.Name(), pair, price)
}
