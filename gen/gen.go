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

func From(dir string) {

	_, err := parse(dir)
	if err != nil {
		logger.Err("gen.from", err)
	}

}
