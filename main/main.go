package main

import (
	"blogger/config"
	"blogger/repo"
	"blogger/webhook"
	"net/url"
)

func main() {

	rep := repo.New("https://github.com/dengzii/RespberryPi", config.Git.AccessToken, "./source")

	wh := webhook.New("0.0.0.0", "/actions/", 8080)

	wh.Register("1", "abcd", func(id string, params url.Values) {
		rep.Remove()
		rep.Clone()
	})

	wh.Listen()
}
