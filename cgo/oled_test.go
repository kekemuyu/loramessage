package cgo

import (
	"loramessage/hzk"
	"testing"
	"time"
)

func TestOLED_Init(t *testing.T) {
	list, _ := hzk.GetMatchedPymb("94664")
	OLED_Text(20, 50, list[1].Pymb, 1)
	OLED_Show()
	for {
		OLED_SCLK_Clr()
		time.Sleep(time.Second * 5)
		OLED_SCLK_Set()
		time.Sleep(time.Second * 5)
	}
}
