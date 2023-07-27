package cgo

//#include"devmem.h"
import "C" //注意：include和import之间不能有空格

const GPIO_BASE = 0x1c20800

func init() {
	VirmapInit(GPIO_BASE) //gpio物理地址映射到sdram
}

func VirmapInit(target int32) {
	C.Openfile(C.long(target))
}

func VirmapDeinit() {
	C.Closefile()
}
