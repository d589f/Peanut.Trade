package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang/internal/api"
	"golang/internal/exchange"
	"golang/internal/marketdata"
	"golang/internal/service"
	"golang/pkg/yellowstone"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	md := marketdata.NewMarketDataService()
	binance := &exchange.BinanceProvider{}
	kucoin := &exchange.KuCoinProvider{}
	uniswap := &exchange.UniswapProvider{}
	cexProviders := []marketdata.ExchangePoller{binance, kucoin, uniswap}
	pairs := []string{"BTC/USDT", "ETH/USDT", "SOL/USDT"} // Можно расширить
	stopCh := make(chan struct{})
	md.PollCEXPrices(cexProviders, pairs, 10*time.Second, stopCh)

	raydium := exchange.NewRaydiumProvider()

	yellowstoneAddr := "solana-yellowstone-grpc.publicnode.com:443"

	yc, err := yellowstone.NewYellowstoneClient(yellowstoneAddr)
	if err != nil {
		log.Fatalf("failed to connect to Yellowstone: %v", err)
	}
	ctx, _ := context.WithCancel(context.Background())
	raydiumUpdate := func(pair string, price float64) { md.UpdatePrice("Raydium", pair, price) }
	raydium.StartYellowstone(ctx, yc, raydiumUpdate, pairs)

	svc := service.NewService(md)
	api.InitHandlers(svc)
	api.RegisterRoutes(app)

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		close(stopCh)
		_ = app.Shutdown()
	}()

	log.Println("Starting Peanut.Trade API on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
