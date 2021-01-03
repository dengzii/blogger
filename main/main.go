package main

import (
	Blogger "blogger"
	"blogger/webhook"
	"fmt"
	"net/url"
)

func main() {

	r := Blogger.GitRepo{
		Url:        "https://github.com/dengzii/RespberryPi",
		StorageDir: "./repo",
	}
	r.Remove()
	//r.Update()
	wh := webhook.New("0.0.0.0", "/actions/", 8080)

	wh.Register("1", "abcd", func(id string, params url.Values) {
		fmt.Println("action triggered")
	})

	wh.Listen()
}
