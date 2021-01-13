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

func SetCreateTime(info os.FileInfo, t time.Time) {
	winFile := info.Sys().(*syscall.Win32FileAttributeData)
	winFile.CreationTime = toFileTime(t)
}

func toFileTime(t time.Time) syscall.Filetime {
	return syscall.Filetime{
		LowDateTime:  uint32(t.UnixNano()),
		HighDateTime: uint32(t.UnixNano() >> 32),
	}
}
