package win32

import (
	"reflect"
	"syscall"
	"unsafe"

	"golang.org/x/exp/errors/fmt"
)

var (
	user32 = syscall.NewLazyDLL("User32.dll")

	findWindowW              = user32.NewProc("FindWindowW")
	enumWindows              = user32.NewProc("EnumWindows")
	enumChildWindows         = user32.NewProc("EnumChildWindows")
	setProcessDPIAware       = user32.NewProc("SetProcessDPIAware")
	getDC                    = user32.NewProc("GetDC")
	releaseDC                = user32.NewProc("ReleaseDC")
	postMessageW             = user32.NewProc("PostMessageW")
	getWindowTextW           = user32.NewProc("GetWindowTextW")
	getWindowTextA           = user32.NewProc("GetWindowTextA")
	getClientRect            = user32.NewProc("GetClientRect")
	setWindowPos             = user32.NewProc("SetWindowPos")
	getWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	getWindowLongPtrW        = user32.NewProc("GetWindowLongPtrW")
	createDesktopW           = user32.NewProc("CreateDesktopW")
)

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
func GetWindowTextW(hwnd uintptr) WideChar {
	str := make([]byte, 200)
	l, _, _ := getWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&str[0])), 200)
	return WideChar(str[:l*2])
}

// return like TheRender
func GetWindowTextA(hwnd uintptr) string {
	str := make([]byte, 200)
	ret, _, _ := getWindowTextA.Call(hwnd, uintptr(unsafe.Pointer(&str[0])), 200)
	return string(str[:ret])
}

// https://learn.microsoft.com/zh-cn/windows/win32/api/winuser/nf-winuser-enumwindows
func EnumWindows(f func(hwnd uintptr, lParam uintptr) uintptr) bool {
	ret, _, _ := enumWindows.Call(syscall.NewCallback(f), 77)
	return ret == 1
}

// type EnumChildWindowsCallbackFunc(hwnd uintptr, lParam uintptr) uintptr

// https://learn.microsoft.com/zh-cn/windows/win32/api/winuser/nf-winuser-enumchildwindows
func EnumChildWindows(hwnd uintptr, f func(hwnd uintptr, lParam uintptr) uintptr) {
	enumChildWindows.Call(hwnd, syscall.NewCallback(f), 77)
}

// https://learn.microsoft.com/zh-cn/windows/win32/api/winuser/nf-winuser-getwindowthreadprocessid
func GetWindowThreadProcessId(hwnd uintptr) uint32 {
	var dw DWORD
	ret, _, _ := getWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&dw)))
	_ = ret // 创建窗口的线程标识符
	return dw
}

// https://learn.microsoft.com/zh-cn/windows/win32/api/psapi/nf-psapi-enumprocesses

type long = int32
type DWORD = uint32
type Hwnd = uintptr

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

const SRCCOPY = 0x00CC0020

const INVALID_HANDLE_VALUE = -1

const (
	TH32CS_SNAPTHREAD = 0x00000004
	TH32CS_SNAPMODULE = 0x00000008
)

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

/*
hwnd HWND，欲定位的窗口句柄
　　hWndInsertAfter HWND，置于hwnd前面的窗口句柄。这个参数必须是窗口的句柄或是下面的值之一：
　  HWND_BOTTOM 将窗口置于其它所有窗口的底部
　　HWND_NOTOPMOST 将窗口置于其它所有窗口的顶部，并位于任何最顶部窗口的后面。如果这个窗口非顶部窗口，这个标记对该窗口并不产生影响
　　HWND_TOP 将窗口置于它所有窗口的顶部
　　HWND_TOPMOST 将窗口置于其它所有窗口的顶部，并位于任何最顶部窗口的前面。即使这个窗口不是活动窗口，也维持最顶部状态
x：

　　int，指定窗口新的X坐标
Y：
　　int，指定窗口新的Y坐标
cx：
　　int，指定窗口新的宽度
cy：
　　int，指定窗口新的高度
wFlags：
　　UINT，指定窗口状态和位置的标记。这个参数使用下面值的组合： SWP_DRAWFRAME 围绕窗口画一个框
　　SWP_FRAMECHANGED 发送一条WM_NCCALCSIZE消息进入窗口，即使窗口的大小没有发生改变。如果不指定这个参数，消息WM_NCCALCSIZE只有在窗口大小发生改变时才发送
　　SWP_HIDEWINDOW 隐藏窗口
　　SWP_NOACTIVATE 不激活窗口
　　SWP_NOCOPYBITS 屏蔽客户区域
　　SWP_NOMOVE 保持当前位置（X和Y参数将被忽略）
　　SWP_NOOWNERZORDER 不改变所有窗口的位置和排列顺序
　　SWP_NOREDRAW 窗口不自动重画
　　SWP_NOREPOSITION 与SWP_NOOWNERZORDER标记相同
　　SWP_NOSENDCHANGING 防止这个窗口接受WM_WINDOWPOSCHANGING消息
　　SWP_NOSIZE 保持当前大小（cx和cy会被忽略）
　　SWP_NOZORDER 保持窗口在列表的当前位置（hWndInsertAfter将被忽略）
　　SWP_SHOWWINDOW 显示窗口
	可以 　　SWP_NOMOVE ||
*/
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowpos
func SetWindowPos(hwnd, hWndInsertAfter uintptr, x, y, cx, cy, uFlags int) uintptr {
	ret, _, _ := setWindowPos.Call(
		hwnd,
		hWndInsertAfter,
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		uintptr(uFlags),
	)
	return ret
}

const (
	GWL_STYLE = -16

	WS_CHILD   = 0x40000000
	WS_VISIBLE = 0x10000000
)

// https://learn.microsoft.com/zh-cn/windows/win32/api/winuser/nf-winuser-getwindowlongptrw
func GetWindowLongPtrW(hwnd uintptr, nIndex int) int {
	ret, _, _ := getWindowLongPtrW.Call(hwnd, uintptr(nIndex))
	return int(ret)
}

const (
	DESKTOP_CREATEMENU      = 0x0004
	DESKTOP_CREATEWINDOW    = 0x0002
	DESKTOP_ENUMERATE       = 0x0040
	DESKTOP_HOOKCONTROL     = 0x0008
	DESKTOP_JOURNALPLAYBACK = 0x0020
	DESKTOP_JOURNALRECORD   = 0x0010
	DESKTOP_READOBJECTS     = 0x0001
	DESKTOP_SWITCHDESKTOP   = 0x0100
	DESKTOP_WRITEOBJECTS    = 0x0080

	GENERIC_ALL = DESKTOP_CREATEMENU |
		DESKTOP_CREATEWINDOW |
		DESKTOP_ENUMERATE |
		DESKTOP_HOOKCONTROL |
		DESKTOP_JOURNALPLAYBACK |
		DESKTOP_JOURNALRECORD |
		DESKTOP_READOBJECTS |
		DESKTOP_SWITCHDESKTOP |
		DESKTOP_WRITEOBJECTS
)

// https://learn.microsoft.com/zh-cn/windows/win32/api/winuser/nf-winuser-createdesktopw
func CreateDesktopW(desktopName string) uintptr {
	ret, _, _ := createDesktopW.Call(
		uintptr(unsafe.Pointer(str(desktopName))),
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(nil)),
		0,
		GENERIC_ALL,
		uintptr(unsafe.Pointer(nil)),
	)
	return ret
}

var (
	wincore = syscall.NewLazyDLL("Api-ms-win-core-version-l1-1-0.dll")

	getFileVersionInfoSizeW = wincore.NewProc("GetFileVersionInfoSizeW")
	getFileVersionInfoW     = wincore.NewProc("GetFileVersionInfoW")
	verQueryValueW          = wincore.NewProc("VerQueryValueW")
)

// https://learn.microsoft.com/zh-cn/windows/win32/api/winver/nf-winver-getfileversioninfosizew
func GetFileVersionInfoSizeW(filename string) uint32 {
	var wh uint32
	size, _, _ := getFileVersionInfoSizeW.Call(
		uintptrStr(filename),
		uintptr(unsafe.Pointer(&wh)),
	)
	return uint32(size)
}

const (
	FILE_VER_GET_LOCALISED  = 0x01
	FILE_VER_GET_NEUTRAL    = 0x02
	FILE_VER_GET_PREFETCHED = 0x04
)

// https://learn.microsoft.com/zh-cn/windows/win32/api/winver/nf-winver-getfileversioninfow
func GetFileVersionInfoW(filename string, size uint32, lpData uintptr) bool {
	var lpd uintptr
	// var bs = make([]byte, size)
	ret, _, _ := getFileVersionInfoW.Call(
		uintptrStr(filename),
		0,
		uintptr(size),
		// lpData,
		// uintptr(unsafe.Pointer(&lpd)),
		lpData,
		// uintptr(unsafe.Pointer(&bs[0])),
	)
	fmt.Println("GetFileVersionInfoW,ret:", ret, "size:", size)
	fmt.Println("lpd:", lpd)
	// fmt.Println("bs:", bs)
	return ret > 0
}

type QueryValueTranslation struct {
	LangID   uint16
	CodePage uint16 // Charset
}

// https://learn.microsoft.com/zh-cn/windows/win32/api/winver/nf-winver-verqueryvaluew
func VerQueryValueW(pBlock uintptr, search string, size uint32) bool {
	// buffer := make([]byte, size)
	var lp uint32
	var buffLen uint32
	ret, _, _ := verQueryValueW.Call(
		pBlock,
		uintptrStr(search),
		// uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&lp)),
		uintptr(unsafe.Pointer(&buffLen)),
	)
	fmt.Println("pBlock", pBlock, "lp:", lp, &lp, uintptr(lp))

	data := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(lp),
		Len:  int(buffLen),
		Cap:  int(buffLen),
	}))
	fmt.Println("data", data)
	// fmt.Println(buffLen, buffer)
	// fmt.Println("get", search, buffLen, ret > 0, "str:", string(buffer), WideChar(buffer).Utf8(), buffer[:buffLen*2], WideChar(buffer[:buffLen]).Utf8())
	return ret > 0
}

// https://learn.microsoft.com/zh-cn/windows/win32/api/winver/nf-winver-verqueryvaluew
func VerQueryValueWTranslation(pBlock uintptr) QueryValueTranslation {
	trans := QueryValueTranslation{}
	// var buff = make([]uint8, 200)
	var buffLen uint32
	verQueryValueW.Call(
		pBlock,
		uintptrStr("\\VarFileInfo\\Translation"),
		uintptr(unsafe.Pointer(&trans)),
		uintptr(unsafe.Pointer(&buffLen)),
	)
	fmt.Println(buffLen, trans, unsafe.Sizeof(trans))
	// fmt.Println(buff[:20])
	return trans
}
