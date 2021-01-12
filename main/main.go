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

	rep := repo.New(config.Git.Repo, config.Git.AccessToken, config.Git.Dir)

	wh := webhook.New("0.0.0.0", "/actions/", 8080)

	wh.Register("blog_push", config.Blog.WebHookAccessToken, func(id string, params url.Values) {
		go func() {
			rep.Remove()
			rep.Clone()
			err := gen.From(config.Git.Dir, &gen.RenderConfig{
				OutputDir:   config.Blog.Dir,
				TemplateDir: config.Git.Dir,
			})
			if err != nil {
				logger.Err("action exec failed", err)
			}
		}()
	})

	wh.Listen()
}
