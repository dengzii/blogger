// +build !windows

package utils

import (
	"os"
	"syscall"
	"time"
)

func GetCreateTime(info os.FileInfo) time.Time {
	stat := info.Sys().(*syscall.Stat_t)
	return time.Unix(stat.Ctim.Sec, 0)
}

func ChangeFileTimeAttr(path string, cTime *time.Time, aTime *time.Time, mTime *time.Time) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	stat := info.Sys().(*syscall.Stat_t)
	if cTime != nil {
		stat.Ctim = syscall.NsecToTimespec(cTime.UnixNano())
	}
	if aTime != nil {
		stat.Atim = syscall.NsecToTimespec(aTime.UnixNano())
	}
	if mTime != nil {
		stat.Mtim = syscall.NsecToTimespec(mTime.UnixNano())
	}
	return syscall.Stat(path, stat)
}
