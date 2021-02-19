package typehelper

import "strconv"

func Int64ToInt(num int64) (res int) {
	res, _ = strconv.Atoi(strconv.FormatInt(num, 10))
	return
}

func MustItoa(x int) (res string) {
	return strconv.Itoa(x)
}

func MustAtoi(str string) (res int) {
	var err error
	if res, err = strconv.Atoi(str); err != nil {
		res = 0
	}
	return
}
