package exchange

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRoundTripper struct {
	resp *http.Response
	err  error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

func TestBinanceProvider_GetLatestPrice_Success(t *testing.T) {
	provider := &BinanceProvider{}
	// Подменяем http.DefaultClient.Transport
	oldTransport := http.DefaultTransport
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(`{"price":"123.45"}`)),
	}
	http.DefaultTransport = &mockRoundTripper{resp: resp}
	defer func() { http.DefaultTransport = oldTransport }()

	price, err := provider.GetLatestPrice("BTC/USDT")
	assert.NoError(t, err)
	assert.Equal(t, 123.45, price)
}

func TestBinanceProvider_GetLatestPrice_Error(t *testing.T) {
	provider := &BinanceProvider{}
	oldTransport := http.DefaultTransport
	resp := &http.Response{StatusCode: 500, Body: ioutil.NopCloser(strings.NewReader(""))}
	http.DefaultTransport = &mockRoundTripper{resp: resp}
	defer func() { http.DefaultTransport = oldTransport }()

	_, err := provider.GetLatestPrice("BTC/USDT")
	assert.Error(t, err)
}
