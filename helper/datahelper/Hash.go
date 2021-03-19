package datahelper

import (
	"crypto/md5"
	"fmt"
	"time"

	"rickonono3/r-blog/helper/typehelper"
)

func MakeHashWithStr(str string) (hashStr string) {
	hash := md5.New()
	hash.Write([]byte(str))
	timeByte := []byte(typehelper.MustItoa64(time.Now().UnixNano()))
	return fmt.Sprintf("%x", hash.Sum(timeByte))
}
