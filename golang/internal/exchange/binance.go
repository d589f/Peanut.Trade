package exchange

import (
	"encoding/json"
	"fmt"
	"golang/internal/marketdata"
	"net/http"
	"strings"
)

type BinanceProvider struct{}

func (b *BinanceProvider) Name() string { return "Binance" }

func (b *BinanceProvider) GetLatestPrice(pair string) (float64, error) {
	symbol := strings.ToUpper(strings.ReplaceAll(pair, "/", "")) // BTC/USDT -> BTCUSDT
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("binance: bad status %d", resp.StatusCode)
	}
	var data struct {
		Price string `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	var price float64
	if _, err := fmt.Sscanf(data.Price, "%f", &price); err != nil {
		return 0, err
	}
	return price, nil
}

func (b *BinanceProvider) Poll(pair string, md marketdata.MarketDataServiceUpdater) {
	price, err := b.GetLatestPrice(pair)
	if err != nil {
		fmt.Printf("[Binance] Error for %s: %v\n", pair, err)
		md.UpdateError(b.Name(), pair, err.Error())
		return
	}
	fmt.Printf("[Binance] Updated price %s: %f\n", pair, price)
	md.UpdatePrice(b.Name(), pair, price)
}
