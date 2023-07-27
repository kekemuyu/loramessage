package cgo
//fix
//#include"key.h"
import "C" //注意：include和import之间不能有空格

func KeyInit() {
	C.key_init()
}

var key_map = map[int]string{
	0b011110001: "HOME",
	0b101110100: "BACK",
	0b011110010: "DIAL",
	0b101111000: "HANGUP",

	0b011110100: "UP",
	0b011111000: "LEFT",
	0b101110001: "RIGHT",
	0b101110010: "DOWN",

	0b110110001: "1",
	0b110110010: "2",
	0b110110100: "3",
	0b110111000: "4",

	0b111010001: "5",
	0b111010010: "6",
	0b111010100: "7",
	0b111011000: "8",

	0b111100001: "9",
	0b111100010: "*",
	0b111100100: "0",
	0b111101000: "#",

	0b111111111: "", //无按键
}

func ReadKey() string {
	return key_map[int(C.read_key())]
}
