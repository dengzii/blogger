package gen

import (
	"blogger/logger"
	"time"
)

type Friend struct {
	Name        string
	Url         string
	Email       string
	Avatar      string
	Description string
}

type Blog struct {
	CategoryArticleMap map[string][]Article
	Category           []string
	Friends            []Friend
	Description        string
	Info               *BlogInfo
}

type BlogInfo struct {
	Title       string
	Description string
	Favicon     string
	Bio         string
}

type Article struct {
	Title        string
	LatestUpdate time.Time
	Category     string

	file *articleFile
}

func From(dir string) *Blog {

	bf, err := parse(dir)
	if err != nil {
		logger.Err("gen.from", err)
	}

	categoryArticles := map[string][]Article{}
	var category []string

	for _, cate := range bf.category {

		var articles []Article
		for _, file := range cate.article {
			articles = append(articles, Article{
				Title:        file.name,
				LatestUpdate: file.modTime,
				Category:     cate.name,
				file:         &file,
			})
		}
		category = append(category, cate.name)
		categoryArticles[cate.name] = articles
	}

	desc, err := bf.description.readString()
	if err != nil {
		return nil
	}
	blogInfo, err := bf.siteInfo.readBlogInfo()
	if err != nil {
		return nil
	}
	friends, err := bf.friend.readFriends()
	if err != nil {
		return nil
	}

	return &Blog{
		CategoryArticleMap: categoryArticles,
		Category:           category,
		Friends:            friends,
		Description:        desc,
		Info:               blogInfo,
	}
}
