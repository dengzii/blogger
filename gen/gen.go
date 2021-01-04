package gen

import (
	"blogger/logger"
	"os"
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
	Category []string
	Articles []Article
	Editing  []Article
	Friend   Friend
}

type Article struct {
	Title        string
	LatestUpdate time.Time
	Category     string
}

func Source(dir string) {

	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
	}

	file, err := os.OpenFile(dir, os.O_RDONLY, os.ModeDir)
	if err != nil {
		logger.E("gen.source", err.Error())
		return
	}
	if file == nil {
		logger.E("gen.source", "cannot open dir "+dir)
		return
	}

}
