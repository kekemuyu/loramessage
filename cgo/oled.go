package cgo

//#include"oled.h"
import "C" //注意：include和import之间不能有空格
import (
	"loramessage/hzk"
)

func init() {
	C.OLED_Init()
}

func OLED_Clear() {
	C.OLED_Clear()
}

func OLED_Unclear() {
	C.OLED_Unclear()
}

func OLED_Show() {
	C.OLED_Show()
}

func OLED_FillRect(x, y, w, h, col byte) {
	C.OLED_FillRect(C.uchar(x), C.uchar(y), C.uchar(w), C.uchar(h), C.uchar(col))
}

func OLED_Text(x, y byte, str string, col byte) {
	if len(str) < 0 {
		return
	}
	for _, v := range str {
		if v <= 128 { //ascii
			if x > 128-8 {
				return
			}
			C.OLED_ShowChar(C.uchar(x), C.uchar(y), C.uchar(byte(v)), C.uchar(col))
			x += 8
		} else { //中文
			if x > 128-16 {
				return
			}
			data := hzk.GetHZKFromUtf8([]byte(string(v)))
			C.OLED_ShowBuffer(C.uchar(x), C.uchar(y), (*C.uchar)(C.CBytes(data)), C.uchar(col))
			x += 16
		}
	}

}

func OLED_ShowChar(x, y, str byte, col byte) {
	C.OLED_ShowChar(C.uchar(x), C.uchar(y), C.uchar(str), C.uchar(col))
}

func OLED_ShowBuffer(x, y byte, data []byte, col byte) {
	C.OLED_ShowBuffer(C.uchar(x), C.uchar(y), (*C.uchar)(C.CBytes(data)), C.uchar(col))
}

func OLED_SCLK_Clr() {
	C.OLED_SCLK_Clr()
	C.OLED_SDIN_Clr()
	C.OLED_RST_Clr()
	C.OLED_DC_Clr()
	C.OLED_CS_Clr()
}

func OLED_SCLK_Set() {
	C.OLED_SCLK_Set()
	C.OLED_SDIN_Set()
	C.OLED_RST_Set()
	C.OLED_DC_Set()
	C.OLED_CS_Set()
}
