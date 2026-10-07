//go:build windows

// Package winutil Windows'ga oid yordamchi funksiyalar:
// administrator tekshiruvi, huquqni ko'tarib qayta ishga tushirish, xabar oynasi.
package winutil

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var (
	shell32           = syscall.NewLazyDLL("shell32.dll")
	user32            = syscall.NewLazyDLL("user32.dll")
	procIsUserAnAdmin = shell32.NewProc("IsUserAnAdmin")
	procShellExecuteW = shell32.NewProc("ShellExecuteW")
	procMessageBoxW   = user32.NewProc("MessageBoxW")
)

// MessageBox turlari
const (
	MBOk         = 0x00000000
	MBIconError  = 0x00000010
	MBIconWarn   = 0x00000030
	MBIconInfo   = 0x00000040
	MBYesNo      = 0x00000004
	MBIconQuest  = 0x00000020
	IDYES        = 6
	swShowNormal = 1
)

// IsAdmin joriy jarayon administrator huquqi bilan ishlayotganini qaytaradi.
func IsAdmin() bool {
	r, _, _ := procIsUserAnAdmin.Call()
	return r != 0
}

// RelaunchElevated dasturni administrator huquqi bilan qayta ishga tushiradi (UAC so'raydi).
// Muvaffaqiyatli bo'lsa joriy jarayon tugatilishi kerak.
func RelaunchElevated(args []string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	verb, _ := syscall.UTF16PtrFromString("runas")
	file, _ := syscall.UTF16PtrFromString(exe)
	params, _ := syscall.UTF16PtrFromString(strings.Join(args, " "))
	dir, _ := syscall.UTF16PtrFromString("")

	r, _, callErr := procShellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		uintptr(unsafe.Pointer(params)),
		uintptr(unsafe.Pointer(dir)),
		uintptr(swShowNormal),
	)
	// ShellExecute 32 dan katta qiymat qaytarsa — muvaffaqiyat
	if r <= 32 {
		if callErr != nil && callErr != syscall.Errno(0) {
			return callErr
		}
		return syscall.Errno(r)
	}
	return nil
}

// MessageBox oddiy Windows xabar oynasini ko'rsatadi.
func MessageBox(title, text string, flags uint) int {
	t, _ := syscall.UTF16PtrFromString(text)
	c, _ := syscall.UTF16PtrFromString(title)
	r, _, _ := procMessageBoxW.Call(0, uintptr(unsafe.Pointer(t)), uintptr(unsafe.Pointer(c)), uintptr(flags))
	return int(r)
}
