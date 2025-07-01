package marketdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarketDataService_UpdateAndGetBestRate(t *testing.T) {
	s := NewMarketDataService()
	s.UpdatePrice("Binance", "BTC/USDT", 2.0)
	s.UpdatePrice("KuCoin", "BTC/USDT", 1.5)
	best, amount := s.GetBestRate(10, "BTC", "USDT")
	assert.Equal(t, "Binance", best)
	assert.Equal(t, 20.0, amount)
}

func TestMarketDataService_UpdateError(t *testing.T) {
	s := NewMarketDataService()
	s.UpdateError("Binance", "BTC/USDT", "fail")
	details := s.GetDetails("BTC", "USDT")
	assert.Equal(t, "fail", details["Binance"].Error)
}
