// Package config LessonMode sozlamalarini (config.json) o'qiydi va yozadi.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Config - dastur sozlamalari. config.json fayli .exe yonida turadi.
type Config struct {
	// Telegram bot tokeni (@BotFather dan). Ustoz kodни shu bot orqali oladi.
	BotToken string `json:"botToken"`
	// Ustozning Telegram ID raqami (masalan @userinfobot dan). Kod shu ID ga boradi.
	AdminID int64 `json:"adminId"`
	// Oynada va xabarlarda ko'rsatiladigan nom.
	SchoolName string `json:"schoolName"`
	// Dars rejimida chalg'ituvchi ilovalarni (Telegram Desktop, o'yinlar...) yopib turish.
	CloseDistractingApps bool `json:"closeDistractingApps"`
	// hosts blokini qayta tekshirish oralig'i (soniya). Standart: 5.
	ReapplySeconds int `json:"reapplySeconds"`
	// Ochish kodining amal qilish muddati (soniya). Standart: 300 (5 daqiqa).
	CodeTTLSeconds int `json:"codeTtlSeconds"`
	// Ketma-ket noto'g'ri urinishlar chegarasi.
	MaxAttempts int `json:"maxAttempts"`
	// Qo'shimcha bloklanadigan domenlar (o'qituvchi qo'shadi).
	ExtraBlocked []string `json:"extraBlocked"`
	// Hech qachon bloklanmaydigan domenlar (oq ro'yxat).
	ExtraAllowed []string `json:"extraAllowed"`
	// Logon paytida avtomatik ishga tushish (--install bilan o'rnatiladi).
	AutoStart bool `json:"autoStart"`
}

// Default standart sozlamalarni qaytaradi.
func Default() Config {
	return Config{
		BotToken:             "",
		AdminID:              0,
		SchoolName:           "LessonMode",
		CloseDistractingApps: true,
		ReapplySeconds:       5,
		CodeTTLSeconds:       300,
		MaxAttempts:          5,
		ExtraBlocked:         []string{},
		ExtraAllowed:         []string{},
		AutoStart:            true,
	}
}

// Path config.json faylining to'liq yo'lini qaytaradi (.exe yonida).
func Path() string {
	exe, err := os.Executable()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(filepath.Dir(exe), "config.json")
}

// Load config.json ni o'qiydi. Fayl bo'lmasa — standart sozlama bilan namuna yozadi.
func Load() (Config, error) {
	c := Default()
	p := Path()
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			// Namuna fayl yozib qo'yamiz, standart bilan davom etamiz.
			_ = Save(c)
			return c, nil
		}
		return c, err
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return c, err
	}
	c.normalize()
	return c, nil
}

func (c *Config) normalize() {
	if c.SchoolName == "" {
		c.SchoolName = "LessonMode"
	}
	if c.ReapplySeconds <= 0 {
		c.ReapplySeconds = 5
	}
	if c.CodeTTLSeconds <= 0 {
		c.CodeTTLSeconds = 300
	}
	if c.MaxAttempts <= 0 {
		c.MaxAttempts = 5
	}
	c.BotToken = strings.TrimSpace(c.BotToken)
}

// Save sozlamalarni config.json ga yozadi.
func Save(c Config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0644)
}

// IsConfigured bot token va admin id to'g'ri kiritilganini tekshiradi.
func (c Config) IsConfigured() bool {
	return c.BotToken != "" && c.AdminID != 0
}
