package crypto

import (
	"testing"
	"time"
)

func TestGetPrice(t *testing.T) {
	cryptoPair := "BTC-USDT"
	priceChan := GetCryptoPairListener(cryptoPair)

	select {
	case price, ok := <-priceChan:
		if !ok {
			t.Fatalf("Channel closed")
		}
		if price.Data[0].InstId != cryptoPair {
			t.Errorf("Expected cryptoPair: %s, received: %s", cryptoPair, price.Data[0].InstId)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("Did not receive price data within the expected time")
	}
}
