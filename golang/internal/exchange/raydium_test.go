package exchange

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRaydiumProvider_GetLatestPrice_NoPrice(t *testing.T) {
	r := NewRaydiumProvider()
	price, err := r.GetLatestPrice("SOL/USDT")
	assert.NoError(t, err)
	assert.Equal(t, 0.0, price)
}

func TestRaydiumProvider_GetLatestPrice_ThreadSafe(t *testing.T) {
	r := NewRaydiumProvider()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.mu.Lock()
			r.prices["SOL/USDT"] = 42.0
			r.mu.Unlock()
			p, err := r.GetLatestPrice("SOL/USDT")
			assert.NoError(t, err)
			assert.Equal(t, 42.0, p)
		}()
	}
	wg.Wait()
}

func TestDecodeRaydiumPoolPrice(t *testing.T) {
	// 80 байт, reserveA=2, reserveB=10
	data := make([]byte, 80)
	data[64] = 2
	data[72] = 10
	price, err := decodeRaydiumPoolPrice(data)
	assert.NoError(t, err)
	assert.Equal(t, 5.0, price)

	// Короткие данные
	_, err = decodeRaydiumPoolPrice([]byte{1, 2, 3})
	assert.Error(t, err)
}
