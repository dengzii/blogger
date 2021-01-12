package utils

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

func Md5Str(str string) string {
	md5h := md5.New()
	_, err := md5h.Write([]byte(str))
	if err != nil {
		return ""
	}
	md5s := fmt.Sprintf("%x", md5h.Sum([]byte("")))
	return md5s
}

func RuneStr(str string) (res string) {
	for _, r := range []rune(str) {
		res += strconv.Itoa(int(r))
	}
	return
}
