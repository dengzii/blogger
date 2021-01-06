package gen

import (
	"errors"
	"html/template"
	"os"
	"strings"
)

type ConvertConfig struct {
	OutputDir   string
	TemplateDir string
}

type Template struct {
	Path      string
	Variables map[string]interface{}
}

type IndexTemplate struct {
	*Template
}

func (that Template) execute(variables interface{}, outputPath string) error {
	t, err := template.ParseFiles(that.Path)
	if err != nil {
		return err
	}
	outputInfo, err := os.Stat(outputPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if outputInfo.IsDir() {
			return errors.New("outputPath must be a file")
		}
		_ = os.Remove(outputPath)
	}

	outputFileName := outputPath[strings.LastIndex(outputPath, string(os.PathSeparator))+1:]
	_ = os.MkdirAll(strings.TrimRight(outputPath, outputFileName), os.ModePerm)

	out, err := os.OpenFile(outputPath, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	if err := t.Execute(out, variables); err != nil {
		return err
	}
	return nil
}

type ArticleTemplate struct {
	*Template
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

func Convert(blog *Blog, config ConvertConfig) {

	if err := config.validate(); err != nil {
		// error
		return
	}

}
