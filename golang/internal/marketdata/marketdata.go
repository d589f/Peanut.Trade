package marketdata

import (
	"fmt"
	"sync"
	"time"
)

type PriceInfo struct {
	Rate       float64
	LastUpdate int64  // unix timestamp
	Error      string // если была ошибка
}

type MarketDataService struct {
	mu     sync.RWMutex
	prices map[string]map[string]*PriceInfo // exchange -> pair -> PriceInfo
}

type MarketDataServiceUpdater interface {
	UpdatePrice(exchangeName, pair string, price float64)
	UpdateError(exchangeName, pair, errMsg string)
}

func NewMarketDataService() *MarketDataService {
	return &MarketDataService{
		prices: make(map[string]map[string]*PriceInfo),
	}
}

func (s *MarketDataService) GetBestRate(amount float64, in, out string) (string, float64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var bestExchange string
	var bestAmount float64
	for exchange, pairs := range s.prices {
		info, ok := pairs[fmt.Sprintf("%s/%s", in, out)]
		if !ok || info.Rate == 0 {
			continue
		}
		output := amount * info.Rate
		if output > bestAmount {
			bestAmount = output
			bestExchange = exchange
		}
	}
	return bestExchange, bestAmount
}

func (s *MarketDataService) GetAllRates(in, out string) map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]float64)
	for exchange, pairs := range s.prices {
		if info, ok := pairs[fmt.Sprintf("%s/%s", in, out)]; ok && info.Rate != 0 {
			result[exchange] = info.Rate
		}
	}
	return result
}

func (s *MarketDataService) UpdatePrice(exchangeName, pair string, price float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.prices[exchangeName] == nil {
		s.prices[exchangeName] = make(map[string]*PriceInfo)
	}
	s.prices[exchangeName][pair] = &PriceInfo{
		Rate:       price,
		LastUpdate: time.Now().Unix(),
		Error:      "",
	}
}

func (s *MarketDataService) UpdateError(exchangeName, pair, errMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.prices[exchangeName] == nil {
		s.prices[exchangeName] = make(map[string]*PriceInfo)
	}
	s.prices[exchangeName][pair] = &PriceInfo{
		Rate:       0,
		LastUpdate: time.Now().Unix(),
		Error:      errMsg,
	}
}

func (s *MarketDataService) GetDetails(in, out string) map[string]*PriceInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	providers := []string{"Binance", "KuCoin", "Uniswap", "Raydium"}
	result := make(map[string]*PriceInfo)
	pair := fmt.Sprintf("%s/%s", in, out)
	for _, ex := range providers {
		info := &PriceInfo{
			Rate:       0,
			LastUpdate: 0,
			Error:      "API not available or no data yet",
		}
		if pairs, ok := s.prices[ex]; ok {
			if realInfo, ok := pairs[pair]; ok {
				info = realInfo
			}
		}
		result[ex] = info
	}
	return result
}

func (s *MarketDataService) PollCEXPrices(providers []ExchangePoller, pairs []string, interval time.Duration, stopCh <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				for _, provider := range providers {
					for _, pair := range pairs {
						provider.Poll(pair, s)
					}
				}
			case <-stopCh:
				return
			}
		}
	}()
}

type ExchangePoller interface {
	Poll(pair string, md MarketDataServiceUpdater)
}
