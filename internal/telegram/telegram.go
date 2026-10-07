// Package telegram ustozning Telegram botiga ochish kodini yuboradi.
package telegram

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// GenerateCode 4 xonali tasodifiy kod yaratadi (xavfsiz tasodif bilan).
func GenerateCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		// zaxira: vaqt asosida
		return fmt.Sprintf("%04d", time.Now().UnixNano()%9000+1000)
	}
	return fmt.Sprintf("%04d", n.Int64()+1000)
}

// SendMessage Telegram bot orqali admin(ustoz)ga matn yuboradi.
func SendMessage(token string, adminID int64, text string) error {
	if strings.TrimSpace(token) == "" || adminID == 0 {
		return errors.New("bot token yoki admin ID sozlanmagan (config.json)")
	}
	api := "https://api.telegram.org/bot" + token + "/sendMessage"
	form := url.Values{}
	form.Set("chat_id", strconv.FormatInt(adminID, 10))
	form.Set("text", text)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.PostForm(api, form)
	if err != nil {
		return fmt.Errorf("Telegram ga ulanib bo'lmadi: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Telegram xatosi (%d): %s", resp.StatusCode, string(body))
	}
	return nil
}

// SendUnlockCode ochish kodini chiroyli xabar bilan ustozga yuboradi.
func SendUnlockCode(token string, adminID int64, code, computer, user string) error {
	text := fmt.Sprintf(
		"🔓 LessonMode — ochish so'rovi\n\n"+
			"Kompyuter: %s\n"+
			"Foydalanuvchi: %s\n"+
			"Vaqt: %s\n\n"+
			"Ochish kodi: %s\n\n"+
			"Bu kodni faqat o'quvchiga dars tugaganda ayting.",
		computer, user, time.Now().Format("15:04:05"), code,
	)
	return SendMessage(token, adminID, text)
}
