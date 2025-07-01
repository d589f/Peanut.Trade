package api

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"golang/internal/marketdata"
	"golang/internal/model"
	"golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	app := fiber.New()
	md := marketdata.NewMarketDataService()
	md.UpdatePrice("Binance", "BTC/USDT", 4.2)
	svc = service.NewService(md)
	RegisterRoutes(app)
	return app
}

func TestEstimateHandler_POST(t *testing.T) {
	app := setupTestApp()
	body := `{"inputAmount":10,"inputCurrency":"BTC","outputCurrency":"USDT"}`
	req := httptest.NewRequest("POST", "/estimate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
	var out model.EstimateResponse
	json.NewDecoder(resp.Body).Decode(&out)
	assert.Equal(t, "Binance", out.ExchangeName)
	assert.Equal(t, 42.0, out.OutputAmount)
}

func TestGetRatesHandler_POST(t *testing.T) {
	app := setupTestApp()
	body := `{"baseCurrency":"BTC","quoteCurrency":"USDT"}`
	req := httptest.NewRequest("POST", "/getRates", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
	var out map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&out)
	assert.Contains(t, out, "rates")
	assert.Contains(t, out, "details")
}

func TestGetRatesHandler_GET(t *testing.T) {
	app := setupTestApp()
	req := httptest.NewRequest("GET", "/getRates?baseCurrency=BTC&quoteCurrency=USDT", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
	var out map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&out)
	assert.Contains(t, out, "rates")
	assert.Contains(t, out, "details")
}

func TestEstimateHandler_GET(t *testing.T) {
	app := setupTestApp()
	req := httptest.NewRequest("GET", "/estimate?inputAmount=10&inputCurrency=BTC&outputCurrency=USDT", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
	var out model.EstimateResponse
	json.NewDecoder(resp.Body).Decode(&out)
	assert.Equal(t, "Binance", out.ExchangeName)
	assert.Equal(t, 42.0, out.OutputAmount)
}
