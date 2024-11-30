package main

import (
	"context"
	"fmt"

	"crypto-monitor/backend"

	"github.com/iaping/go-okx/ws/public"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 设置窗口始终在最前
	runtime.WindowSetAlwaysOnTop(ctx, true)

	go a.subscribeCryptoPrices()
}

func (a *App) subscribeCryptoPrices() {
	cryptoPairs := []string{"BTC-USDT", "ETH-USDT", "XRP-USDT"} // 可以根据需求添加更多币种
	for _, pair := range cryptoPairs {
		tickerChan := backend.GetCryptoPairListener(pair)
		go func(pair string, ch <-chan public.EventTickers) {
			for ticker := range ch {
				priceInfo := map[string]interface{}{
					"pair":  pair,
					"price": ticker.Data[0].Last,
				}
				runtime.EventsEmit(a.ctx, "ticker_update", priceInfo)
				// log.Println("emitting event:", priceInfo)
			}
		}(pair, tickerChan)
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
