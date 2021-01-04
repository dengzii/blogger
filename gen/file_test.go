package gen

import (
	"testing"
)

func TestParse(test *testing.T) {

	blogFile, err := Parse("..\\sample_repo")

	if err != nil {
		test.Fatal(err)
	}

	for _, file := range blogFile.Category {
		test.Log("Name=", file.Name, "Path=", file.Path, len(file.Article))
	}
}

func TestParseEmpty(t *testing.T) {

}
