// Package hosts Windows hosts faylini boshqaradi: bloklangan domenlarni
// maxsus belgilangan blok ichida qo'shadi va olib tashlaydi.
package hosts

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	beginMarker = "# ==== LessonMode BEGIN (avtomatik yaratilgan - qo'lda tahrir qilmang) ===="
	endMarker   = "# ==== LessonMode END ===="
	sinkIP      = "0.0.0.0"
)

// FilePath Windows hosts faylining yo'lini qaytaradi.
func FilePath() string {
	root := os.Getenv("SystemRoot")
	if root == "" {
		root = `C:\Windows`
	}
	return filepath.Join(root, "System32", "drivers", "etc", "hosts")
}

// stripBlock mavjud LessonMode blokini matndan olib tashlaydi.
func stripBlock(content string) string {
	lines := strings.Split(content, "\n")
	var out []string
	inBlock := false
	for _, ln := range lines {
		t := strings.TrimRight(ln, "\r")
		if t == beginMarker {
			inBlock = true
			continue
		}
		if t == endMarker {
			inBlock = false
			continue
		}
		if inBlock {
			continue
		}
		out = append(out, t)
	}
	res := strings.Join(out, "\n")
	// ketma-ket bo'sh qatorlarni biroz tozalaymiz
	for strings.Contains(res, "\n\n\n") {
		res = strings.ReplaceAll(res, "\n\n\n", "\n\n")
	}
	return strings.TrimRight(res, "\n") + "\n"
}

// Apply berilgan domenlarni hosts fayliga bloklab yozadi (administrator huquqi kerak).
func Apply(domains []string) error {
	p := FilePath()
	raw, err := os.ReadFile(p)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	base := stripBlock(string(raw))

	var b strings.Builder
	b.WriteString(base)
	if !strings.HasSuffix(base, "\n") {
		b.WriteString("\n")
	}
	b.WriteString(beginMarker + "\n")
	for _, d := range domains {
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}
		b.WriteString(sinkIP + " " + d + "\n")
		if !strings.HasPrefix(d, "www.") {
			b.WriteString(sinkIP + " www." + d + "\n")
		}
	}
	b.WriteString(endMarker + "\n")

	if err := os.WriteFile(p, []byte(b.String()), 0644); err != nil {
		return err
	}
	FlushDNS()
	return nil
}

// Remove LessonMode blokini hosts faylidan olib tashlaydi.
func Remove() error {
	p := FilePath()
	raw, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	cleaned := stripBlock(string(raw))
	if err := os.WriteFile(p, []byte(cleaned), 0644); err != nil {
		return err
	}
	FlushDNS()
	return nil
}

// IsApplied hosts faylida LessonMode bloki borligini tekshiradi.
func IsApplied() bool {
	p := FilePath()
	f, err := os.Open(p)
	if err != nil {
		return false
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if strings.TrimRight(sc.Text(), "\r") == beginMarker {
			return true
		}
	}
	return false
}

// FlushDNS DNS keshini tozalaydi (o'zgarish darhol kuchga kirishi uchun).
func FlushDNS() {
	cmd := exec.Command("ipconfig", "/flushdns")
	hideWindow(cmd)
	_ = cmd.Run()
}
