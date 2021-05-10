package datahelper

import (
	"crypto/md5"
	"fmt"
	"time"

	"rickonono3/r-blog/helper/typehelper"
)

// MakeHashWithStr 生成带时间戳的哈希字符串
func MakeHashWithStr(str string) (hashStr string) {
	hash := md5.New()
	hash.Write([]byte(str))
	hash.Write([]byte(typehelper.MustItoa64(time.Now().UnixNano())))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
