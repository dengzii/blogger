package repo

import (
	"blogger/logger"
	"blogger/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type Repo interface {
	Update()
	Pull()
	Clone()
}

var Repository *UniversalRepo

type UniversalRepo struct {
	Url           string
	GitDir        string
	DefaultBranch string

	out  *os.File
	repo *git.Repository
	auth *http.TokenAuth
}

type Auth struct {
	http.TokenAuth
}

type GitFileInfo struct {
	Author    string
	Email     string
	CreateAt  *time.Time
	UpdatedAt *time.Time
}

const (
	remoteName = "origin"
	depth      = 1
)

func New(url string, accessToken string, gitDir string) *UniversalRepo {
	var auth *http.TokenAuth
	if len(accessToken) != 0 {
		auth = &http.TokenAuth{
			Token: accessToken,
		}
	}
	Repository = &UniversalRepo{
		Url:           url,
		GitDir:        gitDir,
		DefaultBranch: "master",
		out:           os.Stdout,
		auth:          auth,
	}
	return Repository
}

func (that *UniversalRepo) Update() {

	exist := false
	if that.Exist() {
		exist = true
		r, err := git.PlainOpen(that.GitDir)
		if err != nil {
			checkError(err, "update.open")
			if err == git.ErrRepositoryNotExists {
				exist = false
			} else {
				return
			}
		} else {
			that.repo = r
		}
	}
	if !exist {
		that.Clone()
	}
	that.Pull()
}

func (that *UniversalRepo) Remove() {
	err := os.RemoveAll(that.GitDir)
	checkError(err, "remove")
}

func (that *UniversalRepo) Clone() bool {

	//fs := osfs.New(that.GitDir)
	//storage := filesystem.NewStorage(fs, cache.NewObjectLRUDefault())

	repo, err := git.PlainClone(that.GitDir, false, &git.CloneOptions{
		URL:           that.Url,
		Progress:      that.out,
		Auth:          that.auth,
		RemoteName:    remoteName,
		ReferenceName: plumbing.NewBranchReferenceName(that.DefaultBranch),
		SingleBranch:  true,
		NoCheckout:    false,
		//Depth:         depth,
	})

	if checkError(err, "repo.clone") {
		return false
	}
	logger.D("rep.clone", "done")

	that.repo = repo
	err = that.changeTime("")
	checkError(err, "repo.clone")

	logger.D("repo.ctime", "done")
	return true
}

func (that *UniversalRepo) GetFileInfo(file string) (*GitFileInfo, error) {

	cLog, e := that.repo.Log(&git.LogOptions{
		FileName: &file,
	})
	if e == nil {
		var cTime *time.Time
		var lastCommit *object.Commit
		e = cLog.ForEach(func(commit *object.Commit) error {
			if lastCommit == nil {
				lastCommit = commit
			}
			ct := commit.Author.When
			cTime = &ct
			return nil
		})
		cLog.Close()
		return &GitFileInfo{
			Author:    lastCommit.Author.Name,
			Email:     lastCommit.Author.Email,
			CreateAt:  cTime,
			UpdatedAt: &lastCommit.Author.When,
		}, nil
	}
	return nil, e
}

func (that *UniversalRepo) changeTime(dir string) error {

	fs, err := ioutil.ReadDir(path.Join(that.GitDir, dir))
	if err != nil {
		return err
	}
	for _, info := range fs {
		if info.IsDir() {
			_ = that.changeTime(path.Join(dir, info.Name()))
			continue
		}
		fPath := path.Join(dir, info.Name())
		cLog, e := that.repo.Log(&git.LogOptions{
			FileName: &fPath,
		})
		if e == nil {
			var cTime *time.Time
			var mTime *time.Time
			e = cLog.ForEach(func(commit *object.Commit) error {
				if mTime == nil {
					mt := commit.Author.When
					mTime = &mt
				}
				ct := commit.Author.When
				cTime = &ct
				return nil
			})
			cLog.Close()
			checkError(e, "repo.ctime")
			if cTime != nil && mTime != nil {
				e = utils.ChangeFileTimeAttr(path.Join(that.GitDir, fPath), cTime, nil, mTime)
				checkError(e, "repo.ctime")
			}
		}
	}
	return nil
}

func (that *UniversalRepo) Pull() {

	if that.repo == nil {
		panic("no repository")
	}
	wt, err := that.repo.Worktree()
	if checkError(err, "repo.pull") {
		return
	}
	err = wt.Pull(&git.PullOptions{
		RemoteName:    remoteName,
		ReferenceName: plumbing.NewBranchReferenceName(that.DefaultBranch),
		SingleBranch:  true,
		Depth:         depth,
		Auth:          that.auth,
		Progress:      that.out,
		Force:         false,
	})
	if checkError(err, "pull") {
		return
	}
	logger.D("repo.pull", "done")
}

func (that *UniversalRepo) Exist() bool {
	_, err := os.Stat(that.GitDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		checkError(err, "exist")
	}
	return true
}

func checkError(err error, where string) bool {
	if err != nil {
		logger.E("git."+where, err.Error())
		return true
	}
	return false
}
