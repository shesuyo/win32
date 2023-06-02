package win32

import "syscall"

var (
	kernel32 = syscall.NewLazyDLL("Kernel32.dll")

	getLastError = kernel32.NewProc("GetLastError")
)

func GetLastError() uintptr {
	ret, _, _ := getLastError.Call()
	return ret
}
