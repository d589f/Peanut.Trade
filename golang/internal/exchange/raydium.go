package exchange

import (
	"context"
	"encoding/binary"
	"fmt"
	"golang/pkg/yellowstone"
	"golang/proto"
	"sync"
)

type RaydiumProvider struct {
	mu     sync.RWMutex
	prices map[string]float64 // pair -> price
}

func NewRaydiumProvider() *RaydiumProvider {
	return &RaydiumProvider{
		prices: make(map[string]float64),
	}
}

func (r *RaydiumProvider) Name() string { return "Raydium" }

func (r *RaydiumProvider) GetLatestPrice(pair string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	price, ok := r.prices[pair]
	if !ok {
		return 0, nil // или ошибка, если нужно
	}
	return price, nil
}

func decodeRaydiumPoolPrice(data []byte) (float64, error) {
	if len(data) < 80 {
		return 0, fmt.Errorf("data too short for Raydium pool layout")
	}
	reserveA := binary.LittleEndian.Uint64(data[64:72])
	reserveB := binary.LittleEndian.Uint64(data[72:80])
	if reserveA == 0 {
		return 0, fmt.Errorf("reserveA is zero")
	}
	return float64(reserveB) / float64(reserveA), nil
}

func (r *RaydiumProvider) StartYellowstone(ctx context.Context, yc *yellowstone.YellowstoneClient, updatePrice func(pair string, price float64), pairs []string) {
	// Подписка на Raydium SOL/USDT pool
	req := &proto.SubscribeRequest{
		Accounts: map[string]*proto.SubscribeRequestFilterAccounts{
			"SOLUSDT": {Account: []string{"3nMFwZXwY1s1M5s8vYAHqd4wGs4iSxXE4LRoUMMYqEgF"}},
		},
	}
	go yc.SubscribeAndListen(ctx, req, func(update *proto.SubscribeUpdate) {
		fmt.Println("[Raydium] Yellowstone update received")
		acc := update.GetAccount()
		if acc == nil || acc.GetAccount() == nil {
			fmt.Println("[Raydium] No account data in update")
			return
		}
		data := acc.GetAccount().GetData()
		price, err := decodeRaydiumPoolPrice(data)
		if err != nil {
			fmt.Printf("[Raydium] Error decoding pool price: %v\n", err)
			return
		}
		r.mu.Lock()
		r.prices["SOL/USDT"] = price
		r.mu.Unlock()
		fmt.Printf("[Raydium] Updated price SOL/USDT: %f\n", price)
		updatePrice("SOL/USDT", price)
	})
}
