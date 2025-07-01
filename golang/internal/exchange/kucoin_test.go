package exchange

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKuCoinProvider_GetLatestPrice_Success(t *testing.T) {
	provider := &KuCoinProvider{}
	oldTransport := http.DefaultTransport
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(`{"data":{"price":"321.99"}}`)),
	}
	http.DefaultTransport = &mockRoundTripper{resp: resp}
	defer func() { http.DefaultTransport = oldTransport }()

	price, err := provider.GetLatestPrice("BTC/USDT")
	assert.NoError(t, err)
	assert.Equal(t, 321.99, price)
}

func TestKuCoinProvider_GetLatestPrice_Error(t *testing.T) {
	provider := &KuCoinProvider{}
	oldTransport := http.DefaultTransport
	resp := &http.Response{StatusCode: 500, Body: ioutil.NopCloser(strings.NewReader(""))}
	http.DefaultTransport = &mockRoundTripper{resp: resp}
	defer func() { http.DefaultTransport = oldTransport }()

	_, err := provider.GetLatestPrice("BTC/USDT")
	assert.Error(t, err)
}
