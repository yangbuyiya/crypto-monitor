package main

import (
	"crypto-monitor/pkg/crypto"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/iaping/go-okx/ws/public"
)

type Widget struct {
	Label *widget.Label
	Chan  <-chan public.EventTickers
}

func main() {
	fmt.Println("Crypto Monitor starting...")
	a := app.New()
	w := a.NewWindow("Crypto Monitor")

	// add title
	titleLabel := widget.NewLabel("Crypto Monitor")
	titleLabel.Alignment = fyne.TextAlignCenter

	pairs := []string{"BTC-USDT", "XRP-USDT", "ETH-USDT", "SOL-USDT", "DOGE-USDT"}
	widgets := []Widget{}

	// create labels and channels for each crypto pair
	for _, pair := range pairs {
		ch := crypto.GetCryptoPairListener(pair)
		label := widget.NewLabel(pair + " Price: Loading...")
		label.Alignment = fyne.TextAlignLeading
		widgets = append(widgets, Widget{
			Label: label,
			Chan:  ch,
		})
	}

	content := container.NewVBox(
		titleLabel,
	)

	for _, widget := range widgets {
		content.Add(widget.Label)
	}

	go UpdateUI(a, widgets)

	w.SetContent(content)
	w.Resize(fyne.NewSize(300, 400))
	w.ShowAndRun()
}
