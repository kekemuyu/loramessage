package control

import (
	"github.com/kekemuyu/loramessage/cgo"
)

type AddrBook struct {
	Contact *Slist
	KeyCh   chan string //从control的keych获取按键的管道地址
}

func NewAddrBook(keych chan string) *AddrBook {
	abook := NewList(20, 0, 100, 16, 0, 4, 1)
	abook.Item = append(abook.Item, []string{"xiaoli", "xiaowang"}...)
	return &AddrBook{
		Contact: abook,
		KeyCh:   keych,
	}
}

func (c *AddrBook) Run() {
	contactEdit := NewInput(c.KeyCh)
	c.Contact.Update()
	cgo.OLED_Show()
	for {
		select {
		case key_value := <-c.KeyCh:
			switch key_value {
			case "UP":
				if c.Contact.Selected > 0 {
					c.Contact.Selected -= 1
				}

			case "DOWN":
				if c.Contact.Selected < len(c.Contact.Item) {
					c.Contact.Selected += 1
				}

			case "C":
			case "HOME": //增加联系人

				contactEdit.Run()
				//从编辑框获取输入的联系人名字，添加到联系人列表
				if len(contactEdit.Memo.Text) > 0 {
					c.Contact.Item = append(c.Contact.Item, contactEdit.Memo.Text)

				}
				contactEdit.Init()
				cgo.OLED_Clear()
			case "#":
				contactEdit.Run()
				contactEdit.Init()
				cgo.OLED_Clear()
			case "*":
				return
			}
			c.Contact.Update()
			cgo.OLED_Show()
		default:
		}
	}

}

func (c *AddrBook) View() {

}

func (c *AddrBook) Add() {

}

func (c *AddrBook) Delete() {

}

func (c *AddrBook) Update() {

}
