package Blogger

import (
	"fmt"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"os"
)

type GitRepo struct {
	Url        string
	StorageDir string
	repo       *git.Repository
}

func (that *GitRepo) Update() {
	if that.Exist() {
		that.repo, _ = git.PlainOpen(that.Url)
	} else {
		that.Clone()
	}
}

func (that GitRepo) Remove() {
	err := os.RemoveAll(that.StorageDir)
	checkError(err)
}

func (that GitRepo) Clone() {

	fs := osfs.New(that.StorageDir)
	storage := filesystem.NewStorage(fs, cache.NewObjectLRUDefault())

	_, err := git.Clone(storage, nil, &git.CloneOptions{
		URL:      that.Url,
		Progress: os.Stdout,
	})

	if checkError(err) {
		return
	}
}

func (that GitRepo) Exist() bool {
	_, err := os.Stat(that.StorageDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		checkError(err)
	}
	return true
}

func checkError(err error) bool {
	if err != nil {
		fmt.Println("Git error " + err.Error())
		return true
	}
	return false
}
