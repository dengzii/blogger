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

	wh := webhook.New(config.Blog.Host, "/actions/", config.Blog.Port)

	wh.Register("blog_push", config.Blog.WebHookAccessToken, func(id string, params url.Values) {
		go func() {
			rep.Remove()
			clone := rep.Clone()
			if clone {
				err := gen.From(config.Git.Dir, &gen.RenderConfig{
					OutputDir:   config.Blog.Dir,
					TemplateDir: config.Blog.Template,
				})
				if err != nil {
					logger.Err("main", err)
				}
				logger.D("main", "action done.")
			}
		}()
	})

	wh.Listen()
}
