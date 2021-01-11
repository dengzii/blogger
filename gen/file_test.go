package gen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {

	blogFile, err := parse("..\\sample_repo")

	assert.Nil(t, err)
	assert.NotEmpty(t, blogFile.path)
	assert.NotEmpty(t, blogFile.category)

	assert.Len(t, blogFile.category, 3)
}

func TestParseNotExist(t *testing.T) {

	_, err := parse("..\\not_exist")

	assert.NotNil(t, err)
}

func TestParseReadFile(t *testing.T) {

	blogFile, err := parse("..\\sample_repo")

	assert.Nil(t, err)

	content, err := blogFile.friend.read()

	assert.Nil(t, err)
	assert.NotEmpty(t, content)
}

func TestReadFirstSection(t *testing.T) {

	blogFile, err := parse("..\\sample_repo")
	assert.NotNil(t, blogFile, err)

	art := blogFile.category[0].article[0]
	firstSec, err := art.readFirstSection()
	all, err := art.read()

	assert.NotEmpty(t, firstSec, all)
	assert.NotEqual(t, all, firstSec)
}

func TestFileMd5(t *testing.T) {

	blogFile, _ := parse("..\\sample_repo")
	s, err := blogFile.category[0].article[0].md5()

	assert.Nil(t, err)
	assert.Len(t, s, 32)
}

func TestFileReadBlogInfo(t *testing.T) {

	blogFile, _ := parse("..\\sample_repo")
	if blogFile.validate() != nil {
		t.Log("err")
	}
	bi, err := blogFile.siteInfo.readBlogInfo()

	assert.Nil(t, err)
	assert.NotNil(t, bi)
}
