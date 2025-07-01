package exchange

import (
	"encoding/json"
	"fmt"
	"golang/internal/marketdata"
	"net/http"
	"strings"
)

type KuCoinProvider struct{}

func (k *KuCoinProvider) Name() string { return "KuCoin" }

func (k *KuCoinProvider) GetLatestPrice(pair string) (float64, error) {
	symbol := strings.ToUpper(strings.ReplaceAll(pair, "/", "-")) // BTC/USDT -> BTC-USDT
	url := fmt.Sprintf("https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("kucoin: bad status %d", resp.StatusCode)
	}
	var data struct {
		Data struct {
			Price string `json:"price"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	var price float64
	if _, err := fmt.Sscanf(data.Data.Price, "%f", &price); err != nil {
		return 0, err
	}
	return price, nil
}

func (k *KuCoinProvider) Poll(pair string, md marketdata.MarketDataServiceUpdater) {
	price, err := k.GetLatestPrice(pair)
	if err != nil {
		fmt.Printf("[KuCoin] Error for %s: %v\n", pair, err)
		md.UpdateError(k.Name(), pair, err.Error())
		return
	}
	fmt.Printf("[KuCoin] Updated price %s: %f\n", pair, price)
	md.UpdatePrice(k.Name(), pair, price)
}
