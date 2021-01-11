// +build !linux

package gen

import (
	"os"
	"syscall"
	"time"
)

func getCreateTime(info os.FileInfo) time.Time {
	winFile := info.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(winFile.CreationTime.Nanoseconds()/1e9, 0)
}
