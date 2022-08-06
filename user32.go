package win32

import (
	"syscall"
	"unsafe"
)

var (
	user32 = syscall.NewLazyDLL("User32.dll")

	findWindowW        = user32.NewProc("FindWindowW")
	enumChildWindows   = user32.NewProc("EnumChildWindows")
	setProcessDPIAware = user32.NewProc("SetProcessDPIAware")
	getDC              = user32.NewProc("GetDC")
	releaseDC          = user32.NewProc("ReleaseDC")
	postMessageW       = user32.NewProc("PostMessageW")
	getWindowTextW     = user32.NewProc("GetWindowTextW")
	getWindowTextA     = user32.NewProc("GetWindowTextA")
	getClientRect      = user32.NewProc("GetClientRect")
)

var (
	gdi32 = syscall.NewLazyDLL("Gdi32.dll")

	createCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
	createCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	selectObject           = gdi32.NewProc("SelectObject")
	bitBlt                 = gdi32.NewProc("BitBlt")
	getBitmapBits          = gdi32.NewProc("GetBitmapBits")
	deleteObject           = gdi32.NewProc("DeleteObject")
)

func str(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-findwindoww
func FindWindowW(lpClassName, lpWindowName string) uintptr {
	var arg1 uintptr = 0
	if lpClassName > "" {
		arg1 = uintptr(unsafe.Pointer(str(lpClassName)))
	}
	ret, _, _ := findWindowW.Call(
		arg1,
		uintptr(unsafe.Pointer(str(lpWindowName))),
	)
	return ret
}

// return like T h e R e n d e r
func GetWindowTextW(hwnd uintptr) string {
	str := make([]byte, 200)
	getWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&str[0])), 200)
	return string(str)
}

// return like TheRender
func GetWindowTextA(hwnd uintptr) string {
	str := make([]byte, 200)
	ret, _, _ := getWindowTextA.Call(hwnd, uintptr(unsafe.Pointer(&str[0])), 200)
	return string(str[:ret])
}

// type EnumChildWindowsCallbackFunc(hwnd uintptr, lParam uintptr) uintptr

func EnumChildWindows(hwnd uintptr, f func(hwnd uintptr, lParam uintptr) uintptr) {
	enumChildWindows.Call(hwnd, syscall.NewCallback(f), 200)
}

type long = int32

type Rect struct {
	Left   long
	Top    long
	Right  long
	Bottom long
}

func GetClientRect(hwnd uintptr) Rect {
	rect := Rect{}
	getClientRect.Call(
		hwnd,
		uintptr(unsafe.Pointer(&rect)),
	)
	return rect
}

func GetDC(hwnd uintptr) uintptr {
	ret, _, _ := getDC.Call(hwnd)
	return ret
}

func CreateCompatibleDC(dc uintptr) uintptr {
	cdc, _, _ := createCompatibleDC.Call(dc)
	return cdc
}

func CreateCompatibleBitmap(dc uintptr, width, height long) uintptr {
	bitmap, _, _ := createCompatibleBitmap.Call(
		dc,
		uintptr(width),
		uintptr(height),
	)
	return bitmap
}

func SelectObject(cdc, bitmap uintptr) {
	selectObject.Call(cdc, bitmap)
}

const SRCCOPY = 0x00CC0020

func BitBlt(cdc, dc, action uintptr, width, height long) {
	bitBlt.Call(
		cdc,
		0,
		0,
		uintptr(width),
		uintptr(height),
		dc,
		0,
		0,
		action,
	)
}

func GetBitmapBits(bitmap uintptr, bytesLen int) []byte {
	// 存储顺序为BGRA
	buffer := make([]byte, bytesLen)
	getBitmapBits.Call(
		bitmap,
		uintptr(bytesLen),
		uintptr(unsafe.Pointer(&buffer[0])),
	)
	return buffer
}

func DeleteObject(hwnd uintptr) {
	deleteObject.Call(hwnd)
}

func ReleaseDC(hwnd, dc uintptr) {
	releaseDC.Call(
		hwnd,
		dc,
	)
}

func SetProcessDPIAware() {
	setProcessDPIAware.Call()
}

func PostMessageW(hwnd, key, wparam, lparam uintptr) uintptr {
	ret, _, _ := postMessageW.Call(
		hwnd,
		key,
		wparam,
		lparam,
	)
	return ret
}
