package gen

import (
	"blogger/utils"
	"bufio"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
)

// expressions file type
const (
	typeArticle  = 1
	typeFriends  = 2
	typeSiteInfo = 3
	typeAboutMe  = 4
	typeDir      = 5
	typeCategory = 6
	typeIgnore   = 7
	typeUnknown  = -1
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
	pathSep + ".git" + pathSep,
}

// the pair of string file name and int type
var fileNameTypeMap = map[string]int{}

var typeFileNameMap = map[int]string{}

func init() {
	fileNameTypeMap["*.md"] = typeArticle
	fileNameTypeMap["*.html"] = typeArticle
	fileNameTypeMap["friends.json"] = typeFriends
	fileNameTypeMap["site-info.json"] = typeSiteInfo
	fileNameTypeMap["about-me.md"] = typeAboutMe
	fileNameTypeMap[".ignore"] = typeIgnore

	for name, t := range fileNameTypeMap {
		typeFileNameMap[t] = name
	}

	for i, ele := range excludeFiles {
		s := strings.ReplaceAll(ele, "/", pathSep)
		s = strings.TrimRight(s, pathSep)
		excludeFiles[i] = s
	}
}

type siteFile struct {
	name       string
	fileType   int
	path       string
	modTime    time.Time
	createTime time.Time
}

func (that *siteFile) validate() error {
	s, err := os.Stat(that.path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return errors.New("cannot read directory")
	}
	return nil
}

func (that *siteFile) read() ([]byte, error) {

	if err := that.validate(); err != nil {
		return nil, err
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
func parse(sourceDir string) (bf *blogFile, err error) {

	sourceDir = strings.TrimRight(sourceDir, pathSep)
	i, e := os.Stat(sourceDir)

	if e != nil {
		return nil, e
	}

	sf := toSiteFile(strings.TrimRight(sourceDir, i.Name()), i)
	bf = &blogFile{
		siteFile: &sf,
		category: []categoryFile{},
	}

	err = nil
	dir, e := ioutil.ReadDir(sourceDir)
	if e != nil {
		err = e
		return
	}

	ignoreFileInfo, err := os.Stat(sourceDir + pathSep + typeFileNameMap[typeIgnore])
	if err == nil {
		ignoreFile := toSiteFile(sourceDir, ignoreFileInfo)
		bf.ignore = &ignoreFile
		igStr, err := bf.ignore.readString()
		if err == nil {
			lines := strings.Split(igStr, "\r\n")
			for _, line := range lines {
				line = strings.ReplaceAll(line, "/", pathSep)
				excludeFiles = append(excludeFiles, strings.TrimSpace(line))
			}
		}
	}

	for _, fileInfo := range dir {

		if skipFile(fileInfo) {
			continue
		}

		if fileInfo.IsDir() {
			dirFile := toSiteFile(sourceDir, fileInfo)
			dirFile.fileType = typeCategory
			cateFile := categoryFile{
				siteFile:        &dirFile,
				article:         []articleFile{},
				alternativeName: utils.Md5Str(dirFile.name)[:8],
			}
			articleFileInfos, e := ioutil.ReadDir(dirFile.path)
			if e != nil {
				err = e
				return
			}

			cateArticles := make([]articleFile, 0)
			for _, fi := range articleFileInfos {
				if skipFile(fi) {
					continue
				}
				aSiteFile := toSiteFile(dirFile.path, fi)

				cateArticles = append(cateArticles, articleFile{
					alternativeName: utils.Md5Str(aSiteFile.name[:strings.LastIndex(aSiteFile.name, ".")])[:8],
					siteFile:        &aSiteFile,
				})
			}
			cateFile.article = cateArticles
			bf.category = append(bf.category, cateFile)
		} else {
			f := toSiteFile(sourceDir, fileInfo)
			switch f.fileType {
			case typeArticle:
				// ignore root
			case typeAboutMe:
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
	ignore      *siteFile
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

func (that *friendsFile) readFriends() ([]*Friend, error) {

	blog := Blog{}
	err := that.readJson(&blog)
	if err != nil {
		return []*Friend{}, err
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
	alternativeName string
	*siteFile
}

func (that articleFile) readFirstSection() ([]byte, error) {

	if err := that.validate(); err != nil {
		return nil, err
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
		if len(line) <= 2 {
			lineStr := string(line)
			if lineStr == "\r\n" || lineStr == "\n" {
				return b, nil
			}
		}
		b = append(b, line...)
		if err != nil {
			if err == io.EOF {
				return b, nil
			}
			return b, err
		}
	}
}

type categoryFile struct {
	*siteFile
	article         []articleFile
	alternativeName string
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
	if fileInfo.IsDir() {
		if contains(excludeFiles, pathSep+fileInfo.Name()+pathSep) {
			return true
		}
	} else {
		suffixIndex := strings.LastIndex(fileInfo.Name(), ".")
		if suffixIndex > -1 {
			suffixPattern := "*" + fileInfo.Name()[suffixIndex:]
			if contains(excludeFiles, suffixPattern) {
				return true
			}
		}
		if contains(excludeFiles, fileInfo.Name()) {
			return true
		}
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

func toSiteFile(parentDir string, info os.FileInfo) siteFile {

	path := parentDir + pathSep + info.Name()
	t := typeUnknown

	suffixPattern := info.Name()
	if !info.IsDir() {
		suffixPattern = "*" + suffixPattern[strings.LastIndex(suffixPattern, "."):]
	}

	if info.IsDir() {

		t = typeDir

	} else if fileNameTypeMap[info.Name()] > 0 {

		t = fileNameTypeMap[info.Name()]

	} else if fileNameTypeMap[suffixPattern] > 0 {

		t = fileNameTypeMap[info.Name()]
	}

	return siteFile{
		name:     info.Name(),
		fileType: t,
		path:     path,
		modTime:  info.ModTime(),
	}
}

func getFileCreateTime(info os.FileInfo) time.Time {
	osType := runtime.GOOS
	switch osType {
	case "windows":
		getCreateTime(info)
	case "linux":
		getCreateTime(info)
	}
	return time.Now()
}
