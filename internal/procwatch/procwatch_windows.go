//go:build windows

// Package procwatch dars rejimida chalg'ituvchi ilovalarni yumshoq yopib turadi.
// Faqat ma'lum chalg'ituvchi ilovalar ro'yxati yopiladi; tizim/muhim jarayonlarga tegilmaydi.
package procwatch

import (
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

var (
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	procCreateToolhelp32Snap = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32FirstW      = kernel32.NewProc("Process32FirstW")
	procProcess32NextW       = kernel32.NewProc("Process32NextW")
	procCloseHandle          = kernel32.NewProc("CloseHandle")
)

const th32csSnapProcess = 0x00000002

type processEntry32W struct {
	Size            uint32
	CntUsage        uint32
	ProcessID       uint32
	DefaultHeapID   uintptr
	ModuleID        uint32
	CntThreads      uint32
	ParentProcessID uint32
	PriClassBase    int32
	Flags           uint32
	ExeFile         [260]uint16
}

// DefaultApps - dars paytida yopiladigan chalg'ituvchi ilovalar (image nomlari).
// Bu ro'yxatda FAQAT ijtimoiy/o'yin/chalg'ituvchi ilovalar bor.
var DefaultApps = []string{
	"Telegram.exe", "WhatsApp.exe", "Viber.exe", "Discord.exe", "Skype.exe",
	"Signal.exe", "Line.exe", "KakaoTalk.exe", "WeChat.exe", "Zalo.exe",
	"Steam.exe", "steamwebhelper.exe", "EpicGamesLauncher.exe", "RobloxPlayerBeta.exe",
	"Battle.net.exe", "Spotify.exe", "vlc.exe",
}

// RunningProcesses hozir ishlab turgan jarayonlarning image nomlarini (kichik harfda) qaytaradi.
func RunningProcesses() map[string]bool {
	result := map[string]bool{}
	snap, _, _ := procCreateToolhelp32Snap.Call(uintptr(th32csSnapProcess), 0)
	if snap == uintptr(syscall.InvalidHandle) || snap == 0 {
		return result
	}
	defer procCloseHandle.Call(snap)

	var pe processEntry32W
	pe.Size = uint32(unsafe.Sizeof(pe))
	r, _, _ := procProcess32FirstW.Call(snap, uintptr(unsafe.Pointer(&pe)))
	for r != 0 {
		name := syscall.UTF16ToString(pe.ExeFile[:])
		if name != "" {
			result[strings.ToLower(name)] = true
		}
		r, _, _ = procProcess32NextW.Call(snap, uintptr(unsafe.Pointer(&pe)))
	}
	return result
}

// CloseDistracting ro'yxatdagi ilovalardan ishlab turganlarini yopadi.
func CloseDistracting(apps []string) {
	running := RunningProcesses()
	for _, app := range apps {
		if running[strings.ToLower(app)] {
			cmd := exec.Command("taskkill", "/IM", app, "/F", "/T")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
			_ = cmd.Run()
		}
	}
}
