package control

import (
	"fmt"
	"github.com/kekemuyu/loramessage/cgo"
	"time"
)

type Manager struct {
	KeyCh     chan string //键值管道
	Main_menu *Slist
}

func New() *Manager {
	main_menu := NewList(20, 0, 100, 16, 0, 4, 1)
	main_menu.Item = append(main_menu.Item, []string{"notebook", "message"}...)
	main_menu.Update()
	cgo.OLED_Show()
	KeyEngliash = make(map[string][]string)
	KeyEngliash["0"] = []string{"0"}
	KeyEngliash["1"] = []string{"1"}
	KeyEngliash["2"] = []string{"2", "a", "b", "c"}
	KeyEngliash["3"] = []string{"3", "d", "e", "f"}
	KeyEngliash["4"] = []string{"4", "g", "h", "i"}
	KeyEngliash["5"] = []string{"5", "j", "k", "l"}
	KeyEngliash["6"] = []string{"6", "m", "n", "o"}
	KeyEngliash["7"] = []string{"7", "p", "q", "r", "s"}
	KeyEngliash["8"] = []string{"8", "t", "u", "v"}
	KeyEngliash["9"] = []string{"9", "w", "x", "y", "z"}

	return &Manager{
		make(chan string, 1024),
		main_menu,
	}
}

//主程序运行
//读取按键，执行相应的操作
func (c *Manager) Run() {

	addrbook := NewAddrBook(c.KeyCh)
	message := NewMessage(c.KeyCh) //创建短信息
	go c.keyProcess()              //定时处理按键
	for {
		select {
		case key_value := <-c.KeyCh:
			switch key_value {
			case "UP":
				if c.Main_menu.Selected > 0 {
					c.Main_menu.Selected -= 1
				}

			case "DOWN":
				if c.Main_menu.Selected < len(c.Main_menu.Item) {
					c.Main_menu.Selected += 1
				}

			case "#":
				cgo.OLED_Clear()
				if c.Main_menu.Item[c.Main_menu.Selected] == "notebook" {
					addrbook.Run() //控制权交给addressbook
				} else if c.Main_menu.Item[c.Main_menu.Selected] == "message" {
					message.Run() //控制权交给message
				}

			}
			c.Main_menu.Update()
			cgo.OLED_Show()

		default:
		}
	}
}

func (c *Manager) keyProcess() {
	cgo.KeyInit() //按键gpio
	for {
		time.Sleep(10 * time.Millisecond)
		key_value := cgo.ReadKey()
		if key_value != "" {
			c.KeyCh <- key_value
			fmt.Println(key_value)
		}
	}
}
