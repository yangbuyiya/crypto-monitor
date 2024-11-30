package main

import (
	"fmt"

	"fyne.io/fyne/v2"
)

func UpdateUI(app fyne.App, widgets []Widget) {
	for _, widgetObj := range widgets {
		go func(wg Widget) {
			for ticker := range wg.Chan {
				priceStr := fmt.Sprintf("%s : %s", ticker.Arg.InstId, ticker.Data[0].Last)
				wg.Label.SetText(priceStr)
			}
		}(widgetObj)
	}
}
