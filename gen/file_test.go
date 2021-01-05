package gen

import (
	"testing"
)

func TestParse(test *testing.T) {

	blogFile, err := parse("..\\sample_repo")

	if err != nil {
		test.Fatal(err)
	}

	test.Log(blogFile.path)
	for _, file := range blogFile.category {
		test.Log("name=", file.name, "path=", file.path, len(file.Article))
	}
}

func TestParseNotExist(t *testing.T) {

	_, err := parse("..\\not_exist")

	if err != nil {
		t.Fatal(err)
	}
}

func TestParseReadFile(t *testing.T) {

	blogFile, err := parse("..\\sample_repo")

	if err != nil {
		t.Fatal(err)
	}

	err, content := blogFile.friend.read()

	t.Log(content)
}
