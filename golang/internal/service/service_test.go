package service

import (
	"golang/internal/marketdata"
	"golang/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_Estimate(t *testing.T) {
	md := marketdata.NewMarketDataService()
	md.UpdatePrice("Binance", "BTC/USDT", 4.2)
	s := &Service{MarketData: md}
	resp := s.Estimate(&model.EstimateRequest{InputAmount: 10, InputCurrency: "BTC", OutputCurrency: "USDT"})
	assert.Equal(t, "Binance", resp.ExchangeName)
	assert.Equal(t, 42.0, resp.OutputAmount)
}

func TestService_GetRates(t *testing.T) {
	md := marketdata.NewMarketDataService()
	md.UpdatePrice("Binance", "BTC/USDT", 1.5)
	md.UpdatePrice("KuCoin", "BTC/USDT", 1.2)
	s := &Service{MarketData: md}
	resp := s.GetRates(&model.GetRatesRequest{BaseCurrency: "BTC", QuoteCurrency: "USDT"})
	assert.Len(t, resp.Rates, 2)
	assert.Equal(t, "Binance", resp.Rates[0].ExchangeName)
}
