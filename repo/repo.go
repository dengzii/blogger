package repo

import (
	"blogger/logger"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
)

type Repo interface {
	Update()
	Pull()
	Clone()
}

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

const (
	remoteName = "origin"
	depth      = 1
)

func New(url string, accessToken string, gitDir string) UniversalRepo {
	var auth *http.TokenAuth
	if len(accessToken) != 0 {
		auth = &http.TokenAuth{
			Token: accessToken,
		}
	}
	return UniversalRepo{
		Url:           url,
		GitDir:        gitDir,
		DefaultBranch: "master",
		out:           os.Stdout,
		auth:          auth,
	}
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

	//ref, _ := repo.Head()
	s := "LCD1602_IIC/pi_status.py"
	cLog, e := repo.Log(&git.LogOptions{
		FileName: &s,
	})
	checkError(e, "=")

	cLog.ForEach(func(commit *object.Commit) error {
		logger.D("---", commit.Author.When.Format("2006-01-02 15:01"))
		return nil
	})

	//objs, err := repo.Objects()
	//if err == nil {
	//	_ = objs.ForEach(func(object object.Object) error {
	//		logger.D("repo.clone", object.ID().String())
	//		logger.D("repo.clone", object.Type().String())
	//		if object.Type() == plumbing.BlobObject {
	//			object.Decode()
	//		}
	//		return nil
	//	})
	//}

	if checkError(err, "clone") {
		return false
	}
	that.repo = repo
	logger.D("git.clone", "done")

	return true
}

func (that *UniversalRepo) Pull() {

	if that.repo == nil {
		panic("no repository")
	}
	wt, err := that.repo.Worktree()
	if checkError(err, "pull.wt") {
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
	logger.D("git.pull", "done")
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
