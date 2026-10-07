//go:build !windows

package hosts

import "os/exec"

// hideWindow Windows'dan tashqarida hech narsa qilmaydi (faqat kompilyatsiya uchun).
func hideWindow(cmd *exec.Cmd) {}
