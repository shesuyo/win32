package win32

import "syscall"

func str(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}

type WideChar string

func (wc WideChar) Utf8() string {
	u16 := []uint16{}
	for i := 0; i < len(wc); i += 2 {
		var b uint16
		if i+1 >= len(wc) {
			b = uint16(wc[i])
		} else {
			b = uint16(wc[i+1])*256 + uint16(wc[i])
		}
		u16 = append(u16, b)
	}
	return syscall.UTF16ToString(u16)
}
