package main

import (
	"encoding/json"
	"fmt"
	"loramessage/control"
	"loramessage/db"
	"loramessage/msg"
	"time"
)

func main() {
	ai := autoinc.New(0, 1, 1)
	listMsg := make([]msg.Msg, 0)

	textbs, _ := json.Marshal(control.Ltext{ai.MustID(), 1, "mumu", "hello", true, time.Now()})
	tmsg := msg.Msg{
		Id:      1,
		Datalen: 5,
		Data:    textbs,
	}
	listMsg = append(listMsg, tmsg)
	bs, err := json.Marshal(listMsg)
	if err != nil {
		panic(err)
	}
	db.Update("msg", "msg", bs)
	listMsg2 := make([]msg.Msg, 10)
	bs = db.Read("msg", "msg")
	err = json.Unmarshal(bs, &listMsg2)
	if err != nil {
		panic(err)
	}
	var ttext control.Ltext
	json.Unmarshal(listMsg2[0].Data, &ttext)
	fmt.Println("db.Read:", listMsg2, ttext)
	man := control.New()
	man.Run()
}
