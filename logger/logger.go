package logger

import (
	"fmt"
	time2 "time"
)

func I(tag string, log string) {
	l := fmt.Sprintf("%s I[%s]: %s", time(), tag, log)

	fmt.Println(l)
}

func D(tag string, log string) {
	l := fmt.Sprintf("%s D[%s]: %s", time(), tag, log)

	fmt.Println(l)
}

func E(tag string, log string) {
	l := fmt.Sprintf("%s E[%s]: %s", time(), tag, log)

	fmt.Println(l)
}

func Err(tag string, err error) {
	l := fmt.Sprintf("%s E[%s]: %s", time(), tag, err.Error())

	fmt.Println(l)
}

func time() string {
	return time2.Now().Format("2006/01/02 15:04:05")
}
