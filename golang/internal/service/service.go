package service

import (
	"golang/internal/marketdata"
	"golang/internal/model"
)

type Service struct {
	MarketData *marketdata.MarketDataService
}

func NewService(md *marketdata.MarketDataService) *Service {
	return &Service{MarketData: md}
}

func (s *Service) Estimate(req *model.EstimateRequest) *model.EstimateResponse {
	exchange, amount := s.MarketData.GetBestRate(req.InputAmount, req.InputCurrency, req.OutputCurrency)
	return &model.EstimateResponse{
		ExchangeName: exchange,
		OutputAmount: amount,
	}
}

func (s *Service) GetRates(req *model.GetRatesRequest) *model.GetRatesResponse {
	rates := s.MarketData.GetAllRates(req.BaseCurrency, req.QuoteCurrency)
	result := make([]model.Rate, 0, len(rates))
	for ex, rate := range rates {
		result = append(result, model.Rate{
			ExchangeName: ex,
			Rate:         rate,
		})
	}
	return &model.GetRatesResponse{Rates: result}
}
