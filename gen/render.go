package gen

import (
	"errors"
	"github.com/CloudyKit/jet"
	"io"
	"os"
	"strings"
)

var templateSet *jet.Set

type ConvertConfig struct {
	OutputDir   string
	TemplateDir string
}

type BlogTemplate struct {
	Name      string
	Variables map[string]interface{}
	*jet.Template
}

func (that *BlogTemplate) Render(variables interface{}, outputPath string) error {
	if that.Template == nil {
		t, err := templateSet.GetTemplate(that.Name)
		if err != nil {
			return err
		}
		that.Template = t
	}
	outputInfo, err := os.Stat(outputPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if outputInfo.IsDir() {
			return errors.New("outputPath must be a file: " + outputPath)
		}
		_ = os.Remove(outputPath)
	}

	outputFileName := outputPath[strings.LastIndex(outputPath, string(os.PathSeparator))+1:]
	_ = os.MkdirAll(strings.TrimRight(outputPath, outputFileName), os.ModePerm)

	out, err := os.OpenFile(outputPath, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	if err := that.Execute(out, nil, variables); err != nil {
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

func (that *ConvertConfig) validate() error {

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

	return nil
}

func Render(blog *Blog, config ConvertConfig) error {

	if err := config.validate(); err != nil {
		// error
		return err
	}
	templateSet = jet.NewHTMLSet(config.TemplateDir)

	articleTemplate := &ArticleTemplate{&BlogTemplate{
		Name:      "template_article",
		Variables: nil,
	}}

	var categoryDir,
		categoryName,
		categoryAlternativeName,
		articleOutput string

	var articles []Article
	var allArticle []Article

	for i := range blog.Category {
		categoryName = blog.Category[i]
		categoryAlternativeName = blog.CategoryAlternativeName[i]

		categoryDir = config.OutputDir + pathSep + categoryAlternativeName
		if err := os.Mkdir(categoryDir, os.ModePerm); err != nil {
			if !os.IsExist(err) {
				//logger.E("gen.convert", "mkdir failed:"+categoryAlternativeName)
				return err
			}
		}

		articles = blog.CategoryArticleMap[categoryName]
		for _, a := range articles {
			articleOutput = categoryDir + pathSep + a.AlternativeName + ".html"
			if strings.HasSuffix(a.file.name, ".html") {
				srcF, err := os.OpenFile(a.file.path, os.O_RDONLY, os.ModePerm)
				if err != nil {
					return err
				}
				_ = os.Remove(articleOutput)
				dstF, err := os.OpenFile(articleOutput, os.O_CREATE, os.ModePerm)
				if err != nil {
					return err
				}
				if _, err := io.Copy(dstF, srcF); err != nil {
					return err
				}
			} else {
				a.ReadContent()
				if err := articleTemplate.Render(a, articleOutput); err != nil {
					//logger.E("gen.convert", "generate article failed :"+articleOutput)
					return err
				}
			}
			allArticle = append(allArticle, a)
		}
	}

	indexTemplate := IndexTemplate{&BlogTemplate{
		Name:      "template_index",
		Variables: nil,
	}}

	indexOutput := config.OutputDir + pathSep + "index.html"
	templateVariable := struct {
		Info                *BlogInfo
		Category            []string
		CategoryAlternative []string
		Articles            []Article
	}{
		Info:                blog.Info,
		Category:            blog.Category,
		CategoryAlternative: blog.CategoryAlternativeName,
		Articles:            allArticle,
	}

	if err := indexTemplate.Render(templateVariable, indexOutput); err != nil {
		return err
	}

	return nil
}
