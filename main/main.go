package main

import (
	"blogger/config"
	"blogger/gen"
	"blogger/logger"
	"blogger/repo"
	"blogger/webhook"
	"net/url"
)

func main() {
	blogFile, err := gen.Parse(".\\sample_repo")

	if err != nil {
		panic(err)
	}

	for _, file := range blogFile.Category {
		logger.D("1", file.Name)
		//logger.D("Name=", file.Name, "Path=", file.Path, len(file.Article))
	}
	rep := repo.New("https://github.com/dengzii/RespberryPi", config.Git.AccessToken, "./source")

	wh := webhook.New("0.0.0.0", "/actions/", 8080)

	wh.Register("1", "abcd", func(id string, params url.Values) {
		rep.Remove()
		rep.Clone()
	})

	wh.Listen()
}
