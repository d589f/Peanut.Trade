package exchange

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniswapProvider_GetLatestPrice_UnsupportedPair(t *testing.T) {
	provider := &UniswapProvider{}
	_, err := provider.GetLatestPrice("FAKE/USDT")
	assert.Error(t, err)
}

func TestUniswapProvider_GetLatestPrice_Success(t *testing.T) {
	provider := &UniswapProvider{}
	oldTransport := http.DefaultTransport
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(`{"data":{"pool":{"token0Price":"0.0","token1Price":"42.42"}}}`)),
	}
	http.DefaultTransport = &mockRoundTripper{resp: resp}
	defer func() { http.DefaultTransport = oldTransport }()
	// Добавим тестовый pool
	uniswapPools["ETH/USDC"] = "testpool"
	price, err := provider.GetLatestPrice("ETH/USDC")
	assert.NoError(t, err)
	assert.Equal(t, 42.42, price)
}
