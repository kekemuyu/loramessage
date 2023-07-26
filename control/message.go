package control

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"loramessage/cgo"
	"loramessage/config"
	"loramessage/db"
	"loramessage/msg"
	"loramessage/serial"
	"strconv"
	"strings"
	"time"

	"github.com/issue9/autoinc"
)

//短信息结构
type Ltext struct {
	Id      int64
	Devid   uint32
	Name    string
	Data    string
	Isnew   bool //是否被读过
	Created time.Time
}

type Message struct {
	Com     io.ReadWriteCloser
	Content *Slist
	KeyCh   chan string //从control的keych获取按键的管道地址
	MsgCh   chan msg.Msg
	Autoid  *autoinc.AutoInc
}

func NewMessage(keych chan string) *Message {
	content := NewList(20, 0, 100, 16, 0, 4, 1)
	//从数据库获取消息
	tempbs := db.Read("msg", "msg")
	listMsg := make([]msg.Msg, 0)
	err := json.Unmarshal(tempbs, &listMsg)
	if err != nil {
		panic(err)
	}

	ai := autoinc.New(int64(len(listMsg)), 1, 1)

	if len(listMsg) > 0 {
		for _, v := range listMsg {
			var text Ltext
			err := json.Unmarshal(v.Data, &text)
			if err != nil {
				fmt.Println(err)
				continue
			}
			content.Item = append(content.Item, strconv.Itoa(int(text.Id))+":"+text.Name)
		}

	}

	portnum := config.Cfg.Section("com1").Key("portnum").MustString("")
	baundrate := config.Cfg.Section("com1").Key("baundrate").MustUint()
	com1, err := com.New(portnum, baundrate)
	if err != nil {
		panic(err)
	}
	return &Message{
		Content: content,
		KeyCh:   keych,
		MsgCh:   make(chan msg.Msg, 128),
		Com:     com1,
		Autoid:  ai,
	}
}

func (c *Message) View() {

}

func (c *Message) Add() {

}

func (c *Message) Delete() {

}

func (c *Message) Update() {

}

func (c *Message) Run() {

	c.Content.Update()
	cgo.OLED_Show()
	go c.getMessage()
	for {
		select {
		case key_value := <-c.KeyCh:
			switch key_value {
			case "UP":
				if c.Content.Selected > 0 {
					c.Content.Selected -= 1
				}

			case "DOWN":
				if c.Content.Selected < len(c.Content.Item) {
					c.Content.Selected += 1
				}

			case "C":
			case "DIAL": //发送
				title := c.Content.Item[c.Content.Selected]
				titlesp := strings.Split(title, ":")
				idstr := titlesp[0]
				id, _ := strconv.ParseInt(idstr, 10, 64)
				text, err := ReadTextById(id)
				fmt.Println(text, err)
				if err != nil {
					continue
				}

				textbs, _ := json.Marshal(text)
				msg := msg.Msg{
					Id:      text.Devid,
					Datalen: uint32(len(textbs)),
					Data:    textbs,
				}
				msgbs, _ := json.Marshal(msg)
				fmt.Println(msg, msgbs)
				c.Com.Write(msgbs) //串口发送
			case "#": //确认查看
				title := c.Content.Item[c.Content.Selected]
				titlesp := strings.Split(title, ":")
				idstr := titlesp[0]
				id, _ := strconv.ParseInt(idstr, 10, 64)
				text, err := ReadTextById(id)
				fmt.Println(text, err)
				if err != nil {
					continue
				}

				textStr := "设备号:" + strconv.Itoa(int(text.Devid)) +
					"-" + text.Name + "-" + text.Data + "-" + text.Created.Format("2006-01-02 15:04:05")
				memo := NewMemo("msg", textStr, 0, 0, 127, 63)
				memo.Update()
			case "*":
				return
			}
			c.Content.Update()
			cgo.OLED_Show()
		default:
		}
	}
}

func (c *Message) getMessage() {
	for {
		buf := make([]byte, 8)
		n, err := c.Com.Read(buf)
		if err != nil {
			fmt.Println("c.Com.Read(buf)", err)
			continue
		}
		if n > 0 {
			tempmsg, err := msg.Unpack(buf)
			if err != nil {
				fmt.Println("msg.Unpack(buf)", err)
				continue
			}
			dataBuf := make([]byte, 1024)
			dataSize := 0

			for dataSize < int(tempmsg.Datalen) {
				n, err := c.Com.Read(dataBuf[dataSize:])
				if err != nil {
					fmt.Println("c.Com.Read(dataBuf)", err)
					break
				}
				dataSize += n
			}

			if dataSize >= int(tempmsg.Datalen) {
				tempmsg.Data = dataBuf[:dataSize]
				c.MsgCh <- tempmsg //得到消息
				//消息入库
				tempbs := db.Read("msg", "msg")
				listMsg := make([]msg.Msg, 0)
				err = json.Unmarshal(tempbs, listMsg)
				if err != nil {
					fmt.Println("json.Unmarshal(tempbs,listMsg)", err)
					continue
				}
				listMsg = append(listMsg, tempmsg)
				tepmbs, _ := json.Marshal(listMsg)
				db.Update("message", "message", tepmbs)

				//消息添加到列表
				var text Ltext
				err = json.Unmarshal(tempmsg.Data, &text)
				if err != nil {
					fmt.Println("err=json.Unmarshal(tempmsg.Data,&text)", err)
					continue
				}
				c.Content.Item = append(c.Content.Item, strconv.Itoa(int(text.Id))+":"+text.Name)
			}
		}
	}
}

//根据id获取数据内容
func ReadTextById(id int64) (Ltext, error) {
	bs := db.Read("msg", "msg")
	listMsgs := make([]msg.Msg, 0)
	err := json.Unmarshal(bs, &listMsgs)
	if err != nil {
		panic(err)
	}

	if len(listMsgs) > 0 {
		for _, v := range listMsgs {
			var text Ltext
			err := json.Unmarshal(v.Data, &text)
			if err != nil {
				return text, err
			}
			if text.Id == id {
				return text, nil
			}
		}
	}

	return Ltext{}, errors.New("no text find")
}
