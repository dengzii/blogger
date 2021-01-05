package gen

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// expressions file type
const (
	typeArticle     = 1
	typeFriends     = 2
	typeSiteInfo    = 3
	typeDescription = 4
	typeDir         = 5
	typeCategory    = 6
	typeUnknown     = -1
)

const pathSep = string(os.PathSeparator)

// specify the include files suffix
var includeFiles = []string{
	".md",
	".html",
	".json",
}

// specify the exclude files
var excludeFiles = []string{
	pathSep + ".git",
}

// the pair of string file name and int type
var typeNameMap = map[string]int{}

func init() {
	typeNameMap["*.md"] = typeArticle
	typeNameMap["*.html"] = typeArticle
	typeNameMap["friends.json"] = typeFriends
	typeNameMap["site-info.json"] = typeSiteInfo
	typeNameMap["description.md"] = typeDescription

	for i, ele := range excludeFiles {
		s := strings.ReplaceAll(ele, "/", pathSep)
		s = strings.TrimRight(s, pathSep)
		excludeFiles[i] = s
	}
}

type siteFile struct {
	name     string
	fileType int
	path     string
	modTime  time.Time
}

func (that *siteFile) read() ([]byte, error) {

	s, err := os.Stat(that.path)
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, errors.New("cannot read directory")
	}

	f, err := os.Open(that.path)
	if err != nil {
		return nil, err
	}
	if f != nil {
		defer func() {
			err = f.Close()
			return
		}()
	}

	var b []byte
	bfRd := bufio.NewReader(f)

	for {
		line, err := bfRd.ReadBytes('\n')
		b = append(b, line...)
		if err != nil {
			if err == io.EOF {
				return b, nil
			}
			return b, err
		}
	}
}

func (that *siteFile) readString() (string, error) {

	bytes, err := that.read()

	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (that *siteFile) readJson(v interface{}) error {

	s, err := that.read()
	if err != nil {
		return err
	}

	err = json.Unmarshal(s, &v)
	if err != nil {
		return err
	}
	return nil
}

func (that *siteFile) md5() (string, error) {

	f, err := os.Open(that.path)
	if err != nil {
		return "", err
	}
	if f == nil {
		return "", errors.New("file is nil")
	}

	md5h := md5.New()
	_, err = io.Copy(md5h, f)

	if err != nil {
		return "", err
	}
	md5s := fmt.Sprintf("%x", md5h.Sum([]byte("")))

	return md5s, nil
}

// Check and parse specified dir to blogFile.
func parse(dirPath string) (bf *blogFile, err error) {

	dirPath = strings.TrimRight(dirPath, pathSep)
	i, e := os.Stat(dirPath)

	if e != nil {
		return nil, e
	}

	sf := toSiteFile(strings.TrimRight(dirPath, i.Name()), i)
	bf = &blogFile{
		siteFile: &sf,
		category: []categoryFile{},
	}

	err = nil
	dir, e := ioutil.ReadDir(dirPath)
	if e != nil {
		err = e
		return
	}

	for _, fileInfo := range dir {

		if skipFile(fileInfo) {
			continue
		}

		if fileInfo.IsDir() {
			dirFile := toSiteFile(dirPath, fileInfo)
			dirFile.fileType = typeCategory
			categoryFile := categoryFile{
				siteFile: &dirFile,
				article:  []articleFile{},
			}
			categoryDir, e := ioutil.ReadDir(dirFile.path)
			if e != nil {
				err = e
				return
			}
			for _, fi := range categoryDir {
				if skipFile(fi) {
					continue
				}
				sf := toSiteFile(dirFile.path, fi)
				categoryFile.article = append(categoryFile.article, articleFile{&sf})
			}

			bf.category = append(bf.category, categoryFile)
		} else {
			f := toSiteFile(dirPath, fileInfo)
			switch f.fileType {
			case typeArticle:
				// ignore root
			case typeDescription:
				bf.description = &descriptionFile{&f}
			case typeSiteInfo:
				bf.siteInfo = &siteInfoFile{&f}
			case typeFriends:
				bf.friend = &friendsFile{&f}
			}
		}
	}

	return
}

type blogFile struct {
	friend      *friendsFile
	siteInfo    *siteInfoFile
	description *descriptionFile
	category    []categoryFile
	*siteFile
}

func (that *blogFile) validate() error {
	if that.siteInfo == nil {
		return errors.New("file 'site-info.json' dose not exist")
	}
	return nil
}

type friendsFile struct {
	*siteFile
}

func (that *friendsFile) readFriends() ([]Friend, error) {

	blog := Blog{}
	err := that.readJson(&blog)
	if err != nil {
		return []Friend{}, err
	}
	return blog.Friends, nil
}

type siteInfoFile struct {
	*siteFile
}

func (that siteInfoFile) readBlogInfo() (*BlogInfo, error) {

	info := BlogInfo{}
	err := that.readJson(&info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

type descriptionFile struct {
	*siteFile
}

type articleFile struct {
	*siteFile
}

type categoryFile struct {
	*siteFile
	article []articleFile
}

func contains(slice []string, item ...string) bool {
	for i := range slice {
		for _, contain := range item {
			if slice[i] == contain {
				return true
			}
		}
	}
	return false
}

func skipFile(fileInfo os.FileInfo) bool {

	// skip hidden files
	if strings.HasPrefix(fileInfo.Name(), ".") {
		return true
	}

	// check exclude files
	if contains(excludeFiles, pathSep+fileInfo.Name(), fileInfo.Name()) {
		return true
	}

	// check include files
	for _, include := range includeFiles {
		ignoreCase := strings.ToLower(fileInfo.Name())
		if strings.HasSuffix(ignoreCase, include) {
			return false
		}
	}

	// by default, files that are not included will be excluded
	// directories that are not included will be included
	return !fileInfo.IsDir()
}

func toSiteFile(dirPath string, info os.FileInfo) siteFile {

	path := dirPath + pathSep + info.Name()
	t := typeUnknown

	suffixPattern := info.Name()
	if !info.IsDir() {
		suffixPattern = "*" + suffixPattern[strings.LastIndex(suffixPattern, "."):]
	}

	if info.IsDir() {

		t = typeDir

	} else if typeNameMap[info.Name()] > 0 {

		t = typeNameMap[info.Name()]

	} else if typeNameMap[suffixPattern] > 0 {

		t = typeNameMap[info.Name()]
	}

	return siteFile{
		name:     info.Name(),
		fileType: t,
		path:     path,
		modTime:  info.ModTime(),
	}
}
