package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"crypto-monitor/backend"

	"github.com/iaping/go-okx/ws/public"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx             context.Context
	subscribedPairs map[string]bool
	mu              sync.Mutex
	isCollapsed     bool    // 新增: 控制窗口是否收起
}

// 新增: 窗口收起时的宽度
const (
	CollapsedWidth  = 10
	ExpandedWidth   = 160
	WindowHeight    = 360
)

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		subscribedPairs: make(map[string]bool),
	}
}

// 新增: 初始化窗口位置
func (a *App) initializeWindowPosition() {
	screenWidth, _ := runtime.ScreenGetAll(a.ctx)
	if len(screenWidth) > 0 {
		primaryScreen := screenWidth[0]
		// 将窗口放置在屏幕右侧
		x := primaryScreen.Width - CollapsedWidth
		y := (primaryScreen.Height - WindowHeight) / 2
		
		runtime.WindowSetPosition(a.ctx, x, y)
		runtime.WindowSetSize(a.ctx, CollapsedWidth, WindowHeight)
		a.isCollapsed = true
	}
}

// 新增: 展开窗口
func (a *App) expandWindow() {
	if a.isCollapsed {
		x, y := runtime.WindowGetPosition(a.ctx)
		runtime.WindowSetPosition(a.ctx, x-(ExpandedWidth-CollapsedWidth), y)
		runtime.WindowSetSize(a.ctx, ExpandedWidth, WindowHeight)
		a.isCollapsed = false
	}
}

// 新增: 收起窗口
func (a *App) collapseWindow() {
	if !a.isCollapsed {
		x, y := runtime.WindowGetPosition(a.ctx)
		runtime.WindowSetPosition(a.ctx, x+(ExpandedWidth-CollapsedWidth), y)
		runtime.WindowSetSize(a.ctx, CollapsedWidth, WindowHeight)
		a.isCollapsed = true
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.EventsOn(ctx, "crypto_pairs_changed", a.handleCryptoPairsChanged)
	// 初始化窗口位置
	a.initializeWindowPosition()
	// 监听鼠标进入/离开事件
	runtime.EventsOn(ctx, "mouse_enter", func(...interface{}) { a.expandWindow() })
	runtime.EventsOn(ctx, "mouse_leave", func(...interface{}) { a.collapseWindow() })
}

// handleCryptoPairsChanged 是包装后的事件处理函数
func (a *App) handleCryptoPairsChanged(data ...interface{}) {
	if len(data) < 1 {
		fmt.Println("did not receive any data")
		return
	}

	pairsInterface, ok := data[0].([]interface{})
	if !ok {
		fmt.Println("received data format is incorrect")
		return
	}

	pairs := make([]string, 0, len(pairsInterface))
	for _, pair := range pairsInterface {
		if p, ok := pair.(string); ok {
			pairs = append(pairs, p)
		} else {
			fmt.Println("pair is not a string")
		}
	}

	a.subscribeCryptoPrices(pairs)
}

// subscribeCryptoPrices subscribes to the prices of the given crypto pairs
func (a *App) subscribeCryptoPrices(pairs []string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, pair := range pairs {
		if a.subscribedPairs[pair] {
			continue
		}

		tickerChan, err := backend.GetCryptoPairListener(pair)
		if err != nil {
			priceInfo := map[string]interface{}{
				"pair":  pair,
				"error": err.Error(),
			}
			runtime.EventsEmit(a.ctx, "ticker_subscription_error", priceInfo)
			continue
		}

		a.subscribedPairs[pair] = true

		go func(pair string, ch <-chan public.EventTickers) {
			for ticker := range ch {
				last, _ := strconv.ParseFloat(ticker.Data[0].Last, 64)
				sodUtc0, _ := strconv.ParseFloat(ticker.Data[0].SodUtc0, 64)
				percentage := (last - sodUtc0) / sodUtc0 * 100
				percentageStr := ""
				if percentage > 0 {
					percentageStr = fmt.Sprintf("+%.2f%%", percentage)
				} else {
					percentageStr = fmt.Sprintf("%.2f%%", percentage)
				}
				priceInfo := map[string]interface{}{
					"pair":       pair,
					"price":      ticker.Data[0].Last,
					"percentage": percentageStr,
				}
				runtime.EventsEmit(a.ctx, "ticker_update", priceInfo)
			}

			a.mu.Lock()
			delete(a.subscribedPairs, pair)
			a.mu.Unlock()

			priceInfo := map[string]interface{}{
				"pair": pair,
				"msg":  "ticker subscription closed",
			}
			runtime.EventsEmit(a.ctx, "ticker_subscription_closed", priceInfo)
		}(pair, tickerChan)
	}
}
