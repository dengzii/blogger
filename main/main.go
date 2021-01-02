package main

import (
	"blogger/webhook"
	"fmt"
	"net/url"
)

func main() {

	wh := webhook.New("0.0.0.0", "/actions/", 8080)

	wh.Register("1", "abcd", func(id string, params url.Values) {
		fmt.Println("action triggered")
	})

	wh.Listen()
}
