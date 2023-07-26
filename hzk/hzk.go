package hzk

import (
	"bytes"

	"io/ioutil"
	"os"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const gbkFilePath = "hzk/gbk.font"

// func main() {
// 	list, _ := getMatchedPymb("486") //国
// 	for _, v := range list {
// 		fmt.Println(v)
// 	}
// 	hz, _ := Utf8ToGbk([]byte("中"))
// 	fmt.Println(hz) //得到gbk
// 	//gbk获取字体数据
// 	data := getHZK(hz, "./gbk.font")
// 	fmt.Println(data)
// }

func GetHZKFromUtf8(utf8 []byte) []byte {
	gbk, _ := Utf8ToGbk(utf8)
	return getHZK(gbk, gbkFilePath)
}

func getHZK(gbkcode []byte, hzkFilePath string) []byte {
	re := make([]byte, 32)
	if len(gbkcode) != 2 {
		return re
	}
	hi := int64(gbkcode[0])
	lo := int64(gbkcode[1])
	pos := ((hi-0x81)*191 + (lo - 0x40)) * 32

	fil, err := os.Open(hzkFilePath)
	defer fil.Close()
	if err != nil {
		return re
	}
	fil.Seek(pos, 0)

	fil.Read(re)
	return re
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
