package api

import (
	"golang/internal/model"
	"golang/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var svc *service.Service

func InitHandlers(s *service.Service) {
	svc = s
}

func EstimateHandler(c *fiber.Ctx) error {
	var req model.EstimateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	resp := svc.Estimate(&req)
	return c.JSON(resp)
}

func GetRatesHandler(c *fiber.Ctx) error {
	var req model.GetRatesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	resp := svc.GetRates(&req)
	details := svc.MarketData.GetDetails(req.BaseCurrency, req.QuoteCurrency)
	return c.JSON(fiber.Map{
		"rates":   resp.Rates,
		"details": details,
	})
}

func GetRatesHandlerGET(c *fiber.Ctx) error {
	base := c.Query("baseCurrency")
	quote := c.Query("quoteCurrency")
	if base == "" || quote == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing baseCurrency or quoteCurrency"})
	}
	req := model.GetRatesRequest{
		BaseCurrency:  base,
		QuoteCurrency: quote,
	}
	resp := svc.GetRates(&req)
	details := svc.MarketData.GetDetails(base, quote)
	return c.JSON(fiber.Map{
		"rates":   resp.Rates,
		"details": details,
	})
}

func EstimateHandlerGET(c *fiber.Ctx) error {
	amountStr := c.Query("inputAmount")
	in := c.Query("inputCurrency")
	out := c.Query("outputCurrency")
	if amountStr == "" || in == "" || out == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing inputAmount, inputCurrency или outputCurrency"})
	}
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid inputAmount"})
	}
	req := model.EstimateRequest{
		InputAmount:    amount,
		InputCurrency:  in,
		OutputCurrency: out,
	}
	resp := svc.Estimate(&req)
	return c.JSON(resp)
}

func RegisterRoutes(app *fiber.App) {
	app.Post("/estimate", EstimateHandler)
	app.Post("/getRates", GetRatesHandler)
	app.Get("/getRates", GetRatesHandlerGET)
	app.Get("/estimate", EstimateHandlerGET)
}
