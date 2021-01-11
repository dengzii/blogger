// +build !windows

package gen

import (
	"os"
	"syscall"
	"time"
)

func getCreateTime(info os.FileInfo) time.Time {
	stat := info.Sys().(*syscall.Stat_t)
	return time.Unix(stat.Ctim.Sec, 0)
}
