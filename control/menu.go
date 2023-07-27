package control

import (
	"fmt"
	"github.com/kekemuyu/loramessage/cgo"
	"github.com/kekemuyu/loramessage/hzk"
)

type Slist struct {
	Item       []string
	Selected   int
	Direction  int //0 横向 1 竖向
	X, Y, W, H byte
	Dis_num    int //
}

func NewList(x, y, w, h byte, selected int, dis_num int, dir int) *Slist {
	return &Slist{
		Item:      make([]string, 0),
		Selected:  selected,
		Direction: dir,
		X:         x,
		Y:         y,
		W:         w,
		H:         h,
		Dis_num:   dis_num,
	}
}

func (c *Slist) Update() {
	fmt.Println("c.X, c.Y, c.W, c.H:", c.X, c.Y, c.W, c.H)
	cgo.OLED_FillRect(c.X, c.Y, c.W, c.H, 0)
	selected := c.Selected
	if selected >= len(c.Item) {
		return
	}
	x := c.X
	y := c.Y

	itemCnt := len(c.Item)

	grp := int(itemCnt / c.Dis_num)        //计算页数
	lastPageCnt := itemCnt - c.Dis_num*grp //计算最后一页的元素个数

	selectedPage := int(selected / c.Dis_num)
	selectedLastPageCnt := selected % c.Dis_num

	candis_num := c.Dis_num
	if selectedPage < grp {

	} else {
		candis_num = lastPageCnt
	}

	for i := byte(0); i < byte(candis_num); i++ {
		if c.Direction == 1 { //竖向

			y = c.Y + 16*i
			if i == byte(selectedLastPageCnt) {
				cgo.OLED_FillRect(x, y, 16, 16, 1)
				cgo.OLED_Text(x, y, c.Item[selectedPage*c.Dis_num+int(i)], 0) //选中项反显
			} else {
				cgo.OLED_Text(x, y, c.Item[selectedPage*c.Dis_num+int(i)], 1)
			}
		} else {
			x = c.X + 16*i
			if i == byte(selectedLastPageCnt) {
				cgo.OLED_FillRect(x, y, 16, 16, 1)
				cgo.OLED_Text(x, y, c.Item[selectedPage*c.Dis_num+int(i)], 0) //选中项反显
			} else {
				cgo.OLED_Text(x, y, c.Item[selectedPage*c.Dis_num+int(i)], 1)
			}
		}

	}

}

type Sedit struct {
	X, Y byte
	Name string
	Text string
}

func NewEdit() *Sedit {
	return &Sedit{}
}

func (c *Sedit) Update() {
	cgo.OLED_Clear()
	cgo.OLED_Text(c.X, c.Y, c.Name+":"+c.Text, 1)
}

type Smemo struct {
	Name       string
	Text       string //输入的内容
	X, Y, W, H byte
}

func NewMemo(name, text string, x, y, w, h byte) *Smemo {
	return &Smemo{
		Name: name,
		Text: text,
		X:    x,
		Y:    y,
		W:    w,
		H:    h,
	}
}

func (c *Smemo) Update() {
	cur_x := c.X
	cur_y := c.Y
	cgo.OLED_FillRect(c.X, c.Y, c.W, c.H, 0)

	if len(c.Text) > 0 {
		letter_w := byte(8)
		for _, v := range c.Text {
			cgo.OLED_Text(cur_x, cur_y, string(v), 1)
			if v <= 128 { //ascii
				cur_x += 8
				letter_w = 8
			} else { //中文
				cur_x += 16
				letter_w = 16
			}

			if cur_x-c.X+letter_w > c.W { //换行
				cur_x = c.X
				cur_y += 16
				if (cur_y - c.Y + 16) > c.H {
					return //处理整体上移
				}
			}

		}
	}

}

//输入组件
type Sinput struct {
	Memo           *Smemo
	Preview_hz     *Slist
	Preview_pinyin *Slist
	Pyinput        string //输入的数字
	KeyCh          chan string
	Input_menthod  int //0 :asscii 1:汉字
}

func NewInput(keych chan string) *Sinput {
	memo := NewMemo("contact", "", 0, 0, 128, 48)
	preview_hz := NewList(0, 48, 80, 16, 0, 5, 0)
	preview_pinyin := NewList(80, 48, 48, 16, 0, 1, 0)

	return &Sinput{
		Memo:           memo,
		Preview_hz:     preview_hz,
		Preview_pinyin: preview_pinyin,
		Pyinput:        "",
		KeyCh:          keych,
		Input_menthod:  0,
	}
	return &Sinput{}
}

func (c *Sinput) Init() {
	c.Memo.Text = ""
	c.Preview_hz.Selected = 0
	c.Preview_hz.Item = make([]string, 0)
	c.Preview_pinyin.Selected = 0
	c.Preview_pinyin.Item = make([]string, 0)
	c.Pyinput = ""
	c.Input_menthod = 0
}

func (c *Sinput) Run() {
	for {
		select {
		case key_value := <-c.KeyCh:
			switch key_value {
			case "DIAL": //切换输入法
				if c.Input_menthod == 1 {
					c.Input_menthod = 0
				} else {
					c.Input_menthod = 1
				}

				cgo.OLED_FillRect(0, 48, 127, 16, 0)
				c.Preview_hz.Selected = 0
				c.Preview_pinyin.Selected = 0
				c.Preview_hz.Item = make([]string, 0)
				c.Preview_pinyin.Item = make([]string, 0)
			case "UP": //切换pinyin

				if c.Preview_pinyin.Selected > 0 {
					c.Preview_pinyin.Selected -= 1
				}
				if len(c.Pyinput) > 0 {
					c.Preview_hz.Item = make([]string, 0)
					list, _ := hzk.GetMatchedPymb(c.Pyinput)
					for _, v := range list[c.Preview_pinyin.Selected].Pymb {
						c.Preview_hz.Item = append(c.Preview_hz.Item, string(v))
					}
					c.Preview_hz.Update()
					c.Preview_pinyin.Update()
				}
			case "DOWN": //切换pinyin
				c.Preview_pinyin.Selected += 1
				if c.Preview_pinyin.Selected >= len(c.Preview_pinyin.Item) {
					c.Preview_pinyin.Selected = 0
				}
				if len(c.Pyinput) > 0 {
					c.Preview_hz.Item = make([]string, 0)
					list, _ := hzk.GetMatchedPymb(c.Pyinput)
					for _, v := range list[c.Preview_pinyin.Selected].Pymb {
						c.Preview_hz.Item = append(c.Preview_hz.Item, string(v))
					}
					c.Preview_hz.Update()
					c.Preview_pinyin.Update()
				}

			case "BACK": //退格键
				if len(c.Memo.Text) > 0 {
					strRune := []rune(c.Memo.Text)
					c.Memo.Text = string(strRune[:len(strRune)-1])

				}

			case "LEFT": //切换汉字

				if c.Preview_hz.Selected > 0 {
					c.Preview_hz.Selected -= 1
				}
			case "RIGHT": //切换汉字
				c.Preview_hz.Selected += 1
				if c.Preview_hz.Selected >= len(c.Preview_hz.Item) {
					c.Preview_hz.Selected = 0
				}
			case "#": //确认输入
				if len(c.Preview_hz.Item[c.Preview_hz.Selected]) > 0 {

					c.Memo.Text += c.Preview_hz.Item[c.Preview_hz.Selected]
					c.Pyinput = ""
					c.Preview_hz.Selected = 0
					c.Preview_pinyin.Selected = 0
					c.Preview_hz.Item = make([]string, 0)
					c.Preview_pinyin.Item = make([]string, 0)
				}
			case "*":
				return
			}
			if (c.Input_menthod == 0) && (key_value >= "0") && (key_value <= "9") { //英文处理
				c.Preview_hz.Item = make([]string, 0)
				c.Preview_hz.Item = append(c.Preview_hz.Item, KeyEngliash[key_value]...)

			}
			if (c.Input_menthod == 1) && (key_value >= "0") && (key_value <= "9") { //中文处理
				c.Pyinput += key_value
				fmt.Println(c.Pyinput)
				list, _ := hzk.GetMatchedPymb(c.Pyinput)

				if len(list) > 0 {
					c.Preview_hz.Item = make([]string, 0)
					c.Preview_pinyin.Item = make([]string, 0)
					c.Preview_hz.Selected = 0
					for _, v := range list[0].Pymb {
						c.Preview_hz.Item = append(c.Preview_hz.Item, string(v))
					}

					c.Preview_pinyin.Selected = 0
					for _, v := range list {
						c.Preview_pinyin.Item = append(c.Preview_pinyin.Item, v.Py)
					}
				}
				fmt.Println(c.Preview_hz.Item, c.Preview_pinyin.Item)
			}
			c.Preview_hz.Update()
			c.Preview_pinyin.Update()
			c.Memo.Update() //输入显示到memo组件
			cgo.OLED_Show()
		default:

		}
	}
}

//获取输入字符串
func (c *Sinput) Get() string {
	return c.Memo.Text
}
