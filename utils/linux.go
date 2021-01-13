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

func SetCreateTime(info os.FileInfo, t time.Time) {
	stat := info.Sys().(*syscall.Stat_t)
	stat.Ctim = t.Unix()
}
