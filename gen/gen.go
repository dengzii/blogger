package gen

import (
	"blogger/logger"
	"fmt"
	"os"
	"strings"
)

type Friend struct {
	Name        string
	Url         string
	Email       string
	Avatar      string
	Description string
}

type Category struct {
	Name            string
	AlternativeName string
	Articles        []*Article
}

type About struct {
	Content string
	file    *descriptionFile
}

type Blog struct {
	Category    []*Category
	Friends     []*Friend
	Description *About
	Info        *BlogInfo
}

type BlogInfo struct {
	Title    string
	Keywords string
	Favicon  string
	Bio      string
}

type Article struct {
	Title                   string
	UpdatedAt               string
	CreatedAt               string
	Category                *Category
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
		that.Title, that.UpdatedAt, that.Category.Name, that.AlternativeName, that.AlternativeCategoryName,
	)
}

func From(dir string, renderConfig *RenderConfig) *Blog {

	if err := renderConfig.validate(); err != nil {
		return nil
	}

	bf, err := parse(dir)
	if err != nil {
		logger.Err("gen.from", err)
	}

	var categories []*Category

	for _, cate := range bf.category {

		var articles []*Article
		category := &Category{
			Name:            cate.name,
			AlternativeName: cate.alternativeName,
			Articles:        articles,
		}

		for _, aFile := range cate.article {
			articles = append(articles, fileToArticle(aFile, renderConfig.OutputDir, category))
		}
		category.Articles = articles
		categories = append(categories, category)
	}

	var desc string
	if bf.description != nil {
		desc, err = bf.description.readString()
		if err != nil {
			desc = ""
		}
	}

	blogInfo, err := bf.siteInfo.readBlogInfo()
	if err != nil {
		return nil
	}
	var friends []*Friend
	if bf.friend != nil {
		friends, err = bf.friend.readFriends()
		if err != nil {
			friends = []*Friend{}
		}
	}

	return &Blog{
		Category:    categories,
		Friends:     friends,
		Description: &About{file: bf.description, Content: desc},
		Info:        blogInfo,
	}
}

func fileToArticle(aFile articleFile, outDir string, category *Category) *Article {

	var firstSec string
	if strings.HasSuffix(aFile.name, ".html") {
		firstSec = ""
	} else {
		fSec, err := aFile.readFirstSection()
		if err == nil {
			firstSec = string(fSec)
		}
		firstSec = strings.TrimRight(firstSec, "\r\n")
		firstSec = strings.TrimRight(firstSec, "\n")
	}
	articleName := aFile.name[:strings.LastIndex(aFile.name, ".")]

	out := outDir + category.AlternativeName + pathSep + aFile.alternativeName + ".html"
	outInfo, err := os.Stat(out)

	if err == nil {
		aFile.createTime = getCreateTime(outInfo)
	} else {
		aFile.createTime = aFile.modTime
	}

	return &Article{
		Title:                   articleName,
		UpdatedAt:               aFile.modTime.Format("2006-01-02"),
		CreatedAt:               aFile.createTime.Format("2006-01-02 15:04"),
		Category:                category,
		AlternativeCategoryName: category.AlternativeName,
		AlternativeName:         aFile.alternativeName,
		FirstSection:            firstSec,
		file:                    &aFile,
	}
}
