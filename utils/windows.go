// +build !linux

package utils

import (
	"os"
	"syscall"
	"time"
)

func GetCreateTime(info os.FileInfo) time.Time {
	winFile := info.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(winFile.CreationTime.Nanoseconds()/1e9, 0)
}

func ChangeFileTimeAttr(path string, cTime *time.Time, aTime *time.Time, mTime *time.Time) error {
	pathp, e := syscall.UTF16PtrFromString(path)
	if e != nil {
		return e
	}
	h, e := syscall.CreateFile(pathp,
		syscall.FILE_WRITE_ATTRIBUTES, syscall.FILE_SHARE_WRITE, nil,
		syscall.OPEN_EXISTING, syscall.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if e != nil {
		return e
	}
	defer syscall.Close(h)
	return syscall.SetFileTime(h, toFiletime(cTime), toFiletime(aTime), toFiletime(mTime))
}

func toFiletime(t *time.Time) *syscall.Filetime {
	if t == nil {
		return nil
	}
	a := syscall.NsecToFiletime(t.UnixNano())
	return &a
}
