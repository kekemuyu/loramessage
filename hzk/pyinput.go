package hzk

// "fmt"

//字符串匹配
func strMatch(str1, str2 string) int {
	if len(str1) == 0 || len(str2) == 0 {
		return 0
	}

	n := 0
	ns1 := len(str1)
	ns2 := len(str2)
	ns := 0
	if ns1 <= ns2 {
		ns = ns1
	} else {
		ns = ns2
	}
	for i := 0; i < ns; i++ {

		if str1[i] != str2[i] {

			break
		} else {
			n = i + 1
		}
	}

	if n == len(str1) {
		return 0xFF
	} else {
		return n
	}

}

func GetMatchedPymb(input string) (matchlist []py_index3, re int) {
	mcnt := 0
	bmcnt := 0

	bestmatch := py_index[0]

	for _, v := range py_index {
		n := strMatch(input, v.Py_input)
		if n > 0 {

			if n == 0xFF { //完全匹配
				matchlist = append(matchlist, v)
				mcnt += 1
			} else if n > bmcnt { //查找最佳匹配
				bmcnt = n
				bestmatch = v
			}
		}
	}

	if mcnt == 0 && bmcnt > 0 {
		matchlist = append(matchlist, bestmatch)
		mcnt = bmcnt | 0X80

	}
	re = mcnt
	return
}
