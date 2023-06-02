package win32

import (
	"syscall"
	"unsafe"
)

var (
	gdi32 = syscall.NewLazyDLL("Gdi32.dll")

	createCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
	createCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	selectObject           = gdi32.NewProc("SelectObject")
	bitBlt                 = gdi32.NewProc("BitBlt")
	getBitmapBits          = gdi32.NewProc("GetBitmapBits")
	deleteObject           = gdi32.NewProc("DeleteObject")
	setStretchBltMode      = gdi32.NewProc("SetStretchBltMode")
	getDIBits              = gdi32.NewProc("GetDIBits")
	getObject              = gdi32.NewProc("GetObjectW")
)

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

func SetStretchBltMode(hdcWindow, HALFTONE uintptr) {
	setStretchBltMode.Call(hdcWindow, HALFTONE)
}

// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-getdibits
// https://docs.microsoft.com/en-us/windows/win32/gdi/capturing-an-image
func GetDIBits(hdc, hmb, start, cLines, lpvBits, lpbmi, usage uintptr) uintptr {
	ret, _, _ := getDIBits.Call(
		hdc,
		hmb,
		start,
		cLines,
		lpvBits,
		lpbmi,
		usage,
	)
	return ret
}

func GetObject(h, c, pv uintptr) {
	getObject.Call(
		h,
		c,
		pv,
	)
}
