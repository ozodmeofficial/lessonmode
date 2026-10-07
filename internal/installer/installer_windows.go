//go:build windows

// Package installer LessonMode'ni o'rnatadi va o'chiradi.
// Barcha amallar shaffof va qaytariladigan: oddiy Scheduled Task va fayl nusxasi.
// Administrator istalgan vaqtda --uninstall bilan yoki Task Scheduler orqali olib tashlashi mumkin.
package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const taskName = "LessonMode"

// InstallDir dastur o'rnatiladigan papka.
func InstallDir() string {
	pf := os.Getenv("ProgramFiles")
	if pf == "" {
		pf = `C:\Program Files`
	}
	return filepath.Join(pf, "LessonMode")
}

// InstalledExePath o'rnatilgan .exe yo'li.
func InstalledExePath() string {
	return filepath.Join(InstallDir(), "LessonMode.exe")
}

func hidden(cmd *exec.Cmd) *exec.Cmd {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	return cmd
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// Install dasturni Program Files ga nusxalaydi, config.json ni ko'chiradi va
// logon avtostartини (Scheduled Task) o'rnatadi. Administrator huquqi kerak.
// Qaytaradi: o'rnatilgan exe yo'li.
func Install() (string, error) {
	dir := InstallDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("papka yaratib bo'lmadi: %w", err)
	}

	cur, err := os.Executable()
	if err != nil {
		return "", err
	}
	dstExe := InstalledExePath()

	// Agar allaqachon o'sha joydan ishlayotgan bo'lsak — nusxalashga hojat yo'q.
	if !sameFile(cur, dstExe) {
		if err := copyFile(cur, dstExe); err != nil {
			return "", fmt.Errorf("exe nusxalab bo'lmadi: %w", err)
		}
	}

	// config.json ni yonига ko'chiramiz (agar bor bo'lsa va hali yo'q bo'lsa).
	srcCfg := filepath.Join(filepath.Dir(cur), "config.json")
	dstCfg := filepath.Join(dir, "config.json")
	if _, err := os.Stat(dstCfg); os.IsNotExist(err) {
		if _, err := os.Stat(srcCfg); err == nil {
			_ = copyFile(srcCfg, dstCfg)
		}
	}

	return dstExe, nil
}

// EnableAutoStart logon paytida avtomatik ishga tushadigan Scheduled Task yaratadi.
func EnableAutoStart() error {
	exe := InstalledExePath()
	cmd := hidden(exec.Command("schtasks", "/Create",
		"/TN", taskName,
		"/TR", "\""+exe+"\"",
		"/SC", "ONLOGON",
		"/RL", "HIGHEST",
		"/F",
	))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("avtostart o'rnatilmadi: %v (%s)", err, string(out))
	}
	return nil
}

// DisableAutoStart avtostart vazifasini o'chiradi.
func DisableAutoStart() error {
	cmd := hidden(exec.Command("schtasks", "/Delete", "/TN", taskName, "/F"))
	_ = cmd.Run() // vazifa bo'lmasa ham xato bermaymiz
	return nil
}

func sameFile(a, b string) bool {
	fa, err1 := filepath.Abs(a)
	fb, err2 := filepath.Abs(b)
	if err1 != nil || err2 != nil {
		return false
	}
	return fa == fb
}
