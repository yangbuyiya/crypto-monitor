package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

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
	isLeftSide      bool    // 新增：标记窗口是在左侧还是右侧
	isAtEdge        bool    // 新增：标记是否在屏幕边缘
	isDragging      bool    // 新增：标记是否在拖动状态
}

// 新增: 窗口收起时的宽度
const (
	CollapsedWidth  = 10
	ExpandedWidth   = 160
	WindowHeight    = 360
	EdgeThreshold   = 20    // 新增：边缘检测阈值
	CheckInterval   = 100 // 检查间隔（毫秒）
)

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		subscribedPairs: make(map[string]bool),
	}
}

// 检查窗口是否在边缘
func (a *App) checkEdgePosition() bool {
	// 如果正在拖动，不进行边缘检测
	if a.isDragging {
		return false
	}

	screens, _ := runtime.ScreenGetAll(a.ctx)
	if len(screens) == 0 {
		return false
	}
	
	primaryScreen := screens[0]
	x, _ := runtime.WindowGetPosition(a.ctx)
	
	// 检查是否在左边缘或右边缘
	isAtLeft := x <= EdgeThreshold
	isAtRight := x >= (primaryScreen.Width - EdgeThreshold - CollapsedWidth)
	
	wasAtEdge := a.isAtEdge
	a.isAtEdge = isAtLeft || isAtRight
	a.isLeftSide = isAtLeft

	// 如果边缘状态发生变化，更新窗口状态
	if wasAtEdge != a.isAtEdge {
		if !a.isAtEdge && a.isCollapsed {
			// 如果离开边缘且窗口是收起的，则展开
			a.expandWindow()
		}
	}
	
	return a.isAtEdge
}

// 优化：展开窗口
func (a *App) expandWindow() {
	if !a.isCollapsed {
		return
	}

	_, y := runtime.WindowGetPosition(a.ctx)
	screens, _ := runtime.ScreenGetAll(a.ctx)
	if len(screens) == 0 {
		return
	}
	
	primaryScreen := screens[0]
	
	if a.isLeftSide {
		runtime.WindowSetPosition(a.ctx, 0, y)
	} else {
		runtime.WindowSetPosition(a.ctx, primaryScreen.Width-ExpandedWidth, y)
	}
	
	runtime.WindowSetSize(a.ctx, ExpandedWidth, WindowHeight)
	a.isCollapsed = false
}

// 优化：收起窗口
func (a *App) collapseWindow() {
	if a.isCollapsed || !a.isAtEdge {
		return
	}

	_, y := runtime.WindowGetPosition(a.ctx)
	screens, _ := runtime.ScreenGetAll(a.ctx)
	if len(screens) == 0 {
		return
	}
	
	primaryScreen := screens[0]
	
	if a.isLeftSide {
		runtime.WindowSetPosition(a.ctx, 0, y)
	} else {
		runtime.WindowSetPosition(a.ctx, primaryScreen.Width-CollapsedWidth, y)
	}
	
	runtime.WindowSetSize(a.ctx, CollapsedWidth, WindowHeight)
	a.isCollapsed = true
}

// 新增：设置拖动状态
func (a *App) SetDragging(dragging bool) {
	a.isDragging = dragging
	if dragging {
		// 在拖动时暂停边缘检测
		a.isAtEdge = false
		// 如果窗口是收起的，则展开
		if a.isCollapsed {
			a.expandWindow()
		}
	} else {
		// 恢复边缘检测
		a.checkEdgePosition()
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.EventsOn(ctx, "crypto_pairs_changed", a.handleCryptoPairsChanged)
	runtime.EventsOn(ctx, "set_dragging", func(optionalData ...interface{}) {
		if len(optionalData) > 0 {
			if dragging, ok := optionalData[0].(bool); ok {
				a.SetDragging(dragging)
			}
		}
	})
	
	// 初始化时检查位置
	a.checkEdgePosition()
	
	// 定期检查窗口位置
	go func() {
		ticker := time.NewTicker(time.Millisecond * CheckInterval)
		for range ticker.C {
			a.checkEdgePosition()
		}
	}()
	
	// 监听鼠标事件
	runtime.EventsOn(ctx, "mouse_enter", func(...interface{}) { 
		if a.isAtEdge {
			a.expandWindow()
		}
	})
	runtime.EventsOn(ctx, "mouse_leave", func(...interface{}) { 
		if a.isAtEdge {
			a.collapseWindow()
		}
	})
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
