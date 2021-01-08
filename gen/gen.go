package gen

import (
	"blogger/logger"
	"fmt"
	"strings"
)

type Friend struct {
	Name        string
	Url         string
	Email       string
	Avatar      string
	Description string
}

type Blog struct {
	CategoryArticleMap      map[string][]Article
	Category                []string
	CategoryAlternativeName []string
	Friends                 []Friend
	Description             string
	Info                    *BlogInfo
}

type BlogInfo struct {
	Title       string
	Description string
	Favicon     string
	Bio         string
}

type Article struct {
	Title                   string
	UpdatedAt               string
	CreatedAt               string
	Category                string
	FirstSection            string
	AlternativeName         string
	AlternativeCategoryName string
	Content                 string

	file *articleFile
}

func (that *Article) ReadContent() string {
	content, err := that.file.readString()
	if err != nil {
		logger.E("gen.article.read", err.Error())
	}
	that.Content = content
	return content
}

func (that *Article) String() string {
	return fmt.Sprintf(
		"Article{Title=%s, UpdatedAt=%s, Category=%s, AlternativeName=%s, AlternativeCategoryName=%s}",
		that.Title, that.UpdatedAt, that.Category, that.AlternativeName, that.AlternativeCategoryName,
	)
}

func From(dir string) *Blog {

	bf, err := parse(dir)
	if err != nil {
		logger.Err("gen.from", err)
	}

	categoryArticles := map[string][]Article{}
	var category []string
	var categoryAlternativeName []string

	for _, cate := range bf.category {

		var articles []Article
		for _, aFile := range cate.article {
			var firstSec string
			if strings.HasSuffix(aFile.name, ".html") {
				firstSec = ""
			} else {
				fSec, err := aFile.readFirstSection()
				if err == nil {
					firstSec = string(fSec)
				}
			}
			articleName := aFile.name[:strings.LastIndex(aFile.name, ".")]
			articles = append(articles, Article{
				Title:                   articleName,
				UpdatedAt:               aFile.modTime.Format("2006-01-02"),
				CreatedAt:               aFile.createTime.Format("2006-01-02"),
				Category:                cate.name,
				AlternativeCategoryName: cate.alternativeName,
				AlternativeName:         aFile.alternativeName,
				FirstSection:            firstSec,
				file:                    &aFile,
			})
		}
		category = append(category, cate.name)
		categoryAlternativeName = append(categoryAlternativeName, cate.alternativeName)
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
		CategoryArticleMap:      categoryArticles,
		Category:                category,
		CategoryAlternativeName: categoryAlternativeName,
		Friends:                 friends,
		Description:             desc,
		Info:                    blogInfo,
	}
}
