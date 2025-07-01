package exchange

type ExchangeProvider interface {
	Name() string
	GetLatestPrice(pair string) (float64, error)
}
