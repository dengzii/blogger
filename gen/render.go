package gen

import (
	"blogger/logger"
	"blogger/utils"
	"bytes"
	"errors"
	"github.com/CloudyKit/jet"
	"github.com/yuin/goldmark"
	"io"
	"os"
	"sort"
	"strings"
)

var templateSet *jet.Set

type RenderConfig struct {
	OutputDir   string
	TemplateDir string
}

type BlogTemplate struct {
	Name      string
	Variables map[string]interface{}
	*jet.Template
}

func (that *BlogTemplate) Init() error {
	if that.Template == nil {
		t, err := templateSet.GetTemplate(that.Name)
		if err != nil {
			return err
		}
		that.Template = t
	}
	return nil
}

func (that *BlogTemplate) Render(variables interface{}, outputPath string) error {

	if err := that.Init(); err != nil {
		return err
	}
	outInfo, err := os.Stat(outputPath)
	var fp *os.File
	if err != nil {
		if os.IsNotExist(err) {
			outputFileName := outputPath[strings.LastIndex(outputPath, string(os.PathSeparator))+1:]
			err = os.MkdirAll(strings.TrimRight(outputPath, outputFileName), os.ModePerm)
			if err != nil {
				return err
			}
			fp, err = os.Create(outputPath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if outInfo.IsDir() {
			return errors.New("outputPath must be a file: " + outputPath)
		}
		err = os.Truncate(outputPath, 0)
		if err != nil {
			return err
		}
		fp, err = os.OpenFile(outputPath, os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if err := that.Execute(fp, nil, variables); err != nil {
		return err
	}

	return nil
}

type IndexTemplate struct {
	*BlogTemplate
}

type ArticleTemplate struct {
	*BlogTemplate
}

func (that *RenderConfig) validate() error {

	if len(that.OutputDir) == 0 {
		return errors.New("output dir not specified")
	}
	if len(that.TemplateDir) == 0 {
		return errors.New("template dir not specified")
	}

	inf, err := os.Stat(that.OutputDir)
	if err != nil {
		return err
	}
	if !inf.IsDir() {
		return errors.New("the output dir must be a directory")
	}

	inf, err = os.Stat(that.TemplateDir)
	if err != nil {
		return err
	}
	if !inf.IsDir() {
		return errors.New("the template dir must be a directory")
	}

	if !strings.HasSuffix(that.OutputDir, pathSep) {
		that.OutputDir += pathSep
	}

	if !strings.HasSuffix(that.TemplateDir, pathSep) {
		that.TemplateDir += pathSep
	}

	return nil
}

func Render(blog *Blog, config *RenderConfig) error {

	if err := config.validate(); err != nil {
		// error
		return err
	}
	templateSet = jet.NewHTMLSet(config.TemplateDir)

	articleTemplate := &ArticleTemplate{&BlogTemplate{
		Name:      "template_article",
		Variables: nil,
	}}

	var category *Category
	var categoryDir string

	for i := range blog.Category {

		category = blog.Category[i]
		categoryDir = config.OutputDir + pathSep + category.AlternativeName
		if err := os.Mkdir(categoryDir, os.ModePerm); err != nil {
			if !os.IsExist(err) {
				//logger.E("gen.convert", "mkdir failed:"+categoryAlternativeName)
				return err
			}
		}

		for _, a := range category.Articles {
			if err := renderArticle(blog, articleTemplate, a, config); err != nil {
				return err
			}
		}
	}

	if err := renderIndex(blog, config); err != nil {
		return err
	}

	if err := renderFriends(blog, config); err != nil {
		return err
	}

	if err := renderAbout(blog, config); err != nil {
		return err
	}
	return nil
}

func renderArticle(blog *Blog, template *ArticleTemplate, article *Article, config *RenderConfig) error {

	categoryDir := config.OutputDir + pathSep + article.Category.AlternativeName
	articleOutput := categoryDir + pathSep + article.AlternativeName + ".html"

	if strings.HasSuffix(article.file.name, ".html") {
		srcF, err := os.OpenFile(article.file.path, os.O_RDONLY, os.ModePerm)
		defer func() {
			if srcF != nil {
				srcF.Close()
			}
		}()
		if err != nil {
			return err
		}

		var mod int
		exist, _ := utils.Exist(articleOutput)
		if exist {
			err = os.Truncate(articleOutput, 0)
			if err != nil {
				logger.Err("render.renderArticle", err)
			}
			mod = os.O_WRONLY
		} else {
			mod = os.O_CREATE
		}
		dstF, err := os.OpenFile(articleOutput, mod, os.ModePerm)
		defer func() {
			if dstF != nil {
				dstF.Close()
			}
		}()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstF, srcF); err != nil {
			return err
		}
	} else {
		cnt, err := article.file.read()
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := goldmark.Convert(cnt, &buf); err != nil {
			return err
		}

		var bt []byte
		for true {
			b, e := buf.ReadByte()
			if e != nil {
				if e != io.EOF {
					return e
				} else {
					break
				}
			}
			bt = append(bt, b)
		}

		article.Content = string(bt)
		variables := struct {
			Article *Article
			Info    *BlogInfo
		}{
			Info:    blog.Info,
			Article: article,
		}

		if err := template.Render(variables, articleOutput); err != nil {
			//logger.E("gen.convert", "generate article failed :"+articleOutput)
			return err
		}
	}
	return nil
}

func renderIndex(blog *Blog, config *RenderConfig) error {

	indexTemplate := IndexTemplate{&BlogTemplate{
		Name:      "template_index",
		Variables: nil,
	}}

	var allArticle []*Article

	for _, c := range blog.Category {
		for _, article := range c.Articles {
			allArticle = append(allArticle, article)
		}
	}

	sort.Slice(allArticle, func(i, j int) bool {
		ai := allArticle[i]
		aj := allArticle[j]
		return ai.file.createTime.After(*aj.file.createTime)
	})

	indexOutput := config.OutputDir + pathSep + "index.html"
	var templateVariable = struct {
		Info     *BlogInfo
		Category []*Category
		Articles []*Article
	}{
		Info:     blog.Info,
		Category: blog.Category,
		Articles: allArticle,
	}
	return indexTemplate.Render(templateVariable, indexOutput)
}

func renderFriends(blog *Blog, config *RenderConfig) error {

	friendsTemplate := BlogTemplate{
		Name: "template_friends",
	}
	output := config.OutputDir + pathSep + "friends.html"

	var templateVariable = struct {
		Info    *BlogInfo
		Friends []*Friend
	}{
		Info:    blog.Info,
		Friends: blog.Friends,
	}

	if err := friendsTemplate.Render(templateVariable, output); err != nil {
		return err
	}
	return nil
}

func renderAbout(blog *Blog, config *RenderConfig) error {

	aboutTemplate := BlogTemplate{
		Name: "template_about",
	}
	output := config.OutputDir + pathSep + "about.html"

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(blog.Description.Content), &buf); err != nil {
		return err
	}

	var bt []byte
	for true {
		b, e := buf.ReadByte()
		if e != nil {
			if e != io.EOF {
				return e
			} else {
				break
			}
		}
		bt = append(bt, b)
	}

	var templateVariable = struct {
		Info  *BlogInfo
		About string
	}{
		Info:  blog.Info,
		About: string(bt),
	}

	if err := aboutTemplate.Render(templateVariable, output); err != nil {
		return err
	}
	return nil
}
