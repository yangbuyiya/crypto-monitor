package crypto

import (
	"regexp"

	"github.com/iaping/go-okx/ws/public"
)

func GetCryptoPairListener(cryptoPair string) <-chan public.EventTickers {
	match, err := regexp.MatchString("^[A-Z]+-[A-Z]+$", cryptoPair)
	if err != nil {
		panic(err)
	}
	if !match {
		panic("cryptoPair is not in the correct format. a correct example is BTC-USDT")
	}

	tickerChan := make(chan public.EventTickers)

	handler := func(c public.EventTickers) {
		tickerChan <- c
	}

	handlerError := func(err error) {
		panic(err)
	}

	if err := public.SubscribeTickers(cryptoPair, handler, handlerError, false); err != nil {
		panic(err)
	}

	return tickerChan
}
