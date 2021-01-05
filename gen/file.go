package gen

import (
	"io/ioutil"
	"os"
	"strings"
)

const (
	TypeArticle     = 1
	TypeFriends     = 2
	TypeSiteInfo    = 3
	TypeDescription = 4
	TypeDir         = 5
	TypeCategory    = 6
	TypeUnknown     = -1
)

const pathSep = string(os.PathSeparator)

// specify the include files suffix
var IncludeFiles = []string{
	".md",
	".html",
}

// specify the exclude files
var ExcludeFiles = []string{
	pathSep + ".git",
}

// the pair of string file name and int type
var typeNameMap = map[string]int{}

func init() {
	typeNameMap["*.md"] = TypeArticle
	typeNameMap["*.html"] = TypeArticle
	typeNameMap["friends.md"] = TypeFriends
	typeNameMap["site-info.md"] = TypeSiteInfo
	typeNameMap["description.md"] = TypeDescription

	for i, ele := range ExcludeFiles {
		s := strings.ReplaceAll(ele, "/", pathSep)
		s = strings.TrimRight(s, pathSep)
		ExcludeFiles[i] = s
	}
}

type SiteFile struct {
	Name string
	Type int
	Path string
}

// Check and parse specified dir to BlogFile.
func Parse(dirPath string) (blogFile *BlogFile, err error) {

	dirPath = strings.TrimRight(dirPath, pathSep)

	name := dirPath[strings.LastIndex(dirPath, pathSep)+1:]
	blogFile = &BlogFile{
		Category: []CategoryFile{},
		SiteFile: &SiteFile{
			Name: name,
			Type: TypeDir,
			Path: dirPath,
		},
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
			dirFile.Type = TypeCategory
			categoryFile := CategoryFile{
				SiteFile: &dirFile,
				Article:  []ArticleFile{},
			}
			categoryDir, e := ioutil.ReadDir(dirFile.Path)
			if e != nil {
				err = e
				return
			}
			for _, fi := range categoryDir {
				if skipFile(fi) {
					continue
				}
				sf := toSiteFile(dirFile.Path, fi)
				categoryFile.Article = append(categoryFile.Article, ArticleFile{&sf})
			}

			blogFile.Category = append(blogFile.Category, categoryFile)
		} else {
			f := toSiteFile(dirPath, fileInfo)
			switch f.Type {
			case TypeArticle:
				// ignore root
			case TypeDescription:
				blogFile.Description = DescriptionFile{&f}
			case TypeSiteInfo:
				blogFile.SiteInfo = SiteInfoFile{&f}
			case TypeFriends:
				blogFile.Friend = FriendsFile{&f}
			}
		}
	}

	return
}

type BlogFile struct {
	Friend      FriendsFile
	SiteInfo    SiteInfoFile
	Description DescriptionFile
	Category    []CategoryFile
	*SiteFile
}

type FriendsFile struct {
	*SiteFile
}

type SiteInfoFile struct {
	*SiteFile
}

type DescriptionFile struct {
	*SiteFile
}

type ArticleFile struct {
	*SiteFile
}

type CategoryFile struct {
	*SiteFile
	Article []ArticleFile
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
	if contains(ExcludeFiles, pathSep+fileInfo.Name(), fileInfo.Name()) {
		return true
	}

	// check include files
	for _, include := range IncludeFiles {
		ignoreCase := strings.ToLower(fileInfo.Name())
		if strings.HasSuffix(ignoreCase, include) {
			return false
		}
	}

	// by default, files that are not included will be excluded
	// directories that are not included will be included
	return !fileInfo.IsDir()
}

func toSiteFile(dirPath string, info os.FileInfo) SiteFile {

	path := dirPath + pathSep + info.Name()
	t := TypeUnknown

	suffixPattern := info.Name()
	if !info.IsDir() {
		suffixPattern = "*" + suffixPattern[strings.LastIndex(suffixPattern, "."):]
	}

	if info.IsDir() {

		t = TypeDir

	} else if typeNameMap[info.Name()] > 0 {

		t = typeNameMap[info.Name()]

	} else if typeNameMap[suffixPattern] > 0 {

		t = typeNameMap[info.Name()]
	}

	return SiteFile{
		Name: info.Name(),
		Type: t,
		Path: path,
	}
}
