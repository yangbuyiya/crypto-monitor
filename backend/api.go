package backend

import (
	"fmt"
	"regexp"

	"github.com/iaping/go-okx/ws/public"
)

func GetCryptoPairListener(cryptoPair string) (<-chan public.EventTickers, error) {
	match, err := regexp.MatchString("^[A-Z]+-[A-Z]+$", cryptoPair)
	if err != nil {
		return nil, fmt.Errorf("cryptoPair format error: %w", err)
	}
	if !match {
		return nil, fmt.Errorf("cryptoPair format error, correct example is BTC-USDT")
	}

	tickerChan := make(chan public.EventTickers)

	handler := func(c public.EventTickers) {
		tickerChan <- c
	}

	handlerError := func(err error) {
		fmt.Printf("SubscribeTickers error: %v\n", err)
		close(tickerChan)
	}

	if err := public.SubscribeTickers(cryptoPair, handler, handlerError, false); err != nil {
		return nil, fmt.Errorf("SubscribeTickers failed: %w", err)
	}

	return tickerChan, nil
}
