//go:build !windows

package procwatch

// DefaultApps - Windows'dan tashqarida bo'sh (faqat kompilyatsiya uchun).
var DefaultApps = []string{}

// CloseDistracting Windows'dan tashqarida hech narsa qilmaydi.
func CloseDistracting(apps []string) {}

// RunningProcesses Windows'dan tashqarida bo'sh ro'yxat qaytaradi.
func RunningProcesses() map[string]bool { return map[string]bool{} }
