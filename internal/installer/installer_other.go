//go:build !windows

package installer

// Windows'dan tashqarida stub (faqat kompilyatsiya uchun).
func InstallDir() string       { return "" }
func InstalledExePath() string { return "" }
func Install() (string, error) { return "", nil }
func EnableAutoStart() error   { return nil }
func DisableAutoStart() error  { return nil }
