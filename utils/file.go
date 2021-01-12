package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

const pathSep = string(os.PathSeparator)

func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CopyDir(src string, dest string) error {
	srcOriginal := src
	err := filepath.Walk(src, func(src string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			dstF := strings.Replace(src, srcOriginal, dest, -1)
			e := CopyFile(src, dstF)
			if e != nil {
				return e
			}
		}
		return nil
	})
	return err
}

func CopyFile(src, dst string) error {

	dstSlices := strings.Split(dst, pathSep)
	dstDir := ""
	for i := 0; i < len(dstSlices)-1; i++ {
		dstDir = dstDir + dstSlices[i] + pathSep
	}
	exist, err := Exist(dstDir)
	if !exist {
		err := os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	srcF, err := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	defer func() {
		if srcF != nil {
			srcF.Close()
		}
	}()
	if err != nil {
		return err
	}

	var mod int
	exist, _ = Exist(dst)
	if exist {
		err = os.Truncate(dst, 0)
		if err != nil {
			return err
		}
		mod = os.O_WRONLY
	} else {
		mod = os.O_CREATE
	}
	dstF, err := os.OpenFile(dst, mod, os.ModePerm)
	defer func() {
		if dstF != nil {
			dstF.Close()
		}
	}()
	if err != nil {
		return err
	}

	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}
	return nil
}
