//go:build windows

// LessonMode — sinf uchun diqqatni jamlash (lockdown) dasturi.
//
// QONUNIY FOYDALANISH: faqat o'zingiz boshqarishga vakolatli kompyuterlarda
// (maktab/sinf kompyuterlari) va o'quvchilar xabardor bo'lgan holda ishlating.
// Dastur SHAFFOF: doim ustda turuvchi oyna ko'rinib turadi, administrator uni
// istalgan vaqtda --uninstall bilan yoki hosts faylini tahrirlab olib tashlashi mumkin.
// Dastur tizim funksiyalarini (qayta yuklash, Task Manager) bloklamaydi va
// jarayonni o'ldirishga qarshilik qilmaydi.
package main

import (
	_ "embed"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"lessonmode/internal/blocklist"
	"lessonmode/internal/config"
	"lessonmode/internal/hosts"
	"lessonmode/internal/installer"
	"lessonmode/internal/procwatch"
	"lessonmode/internal/telegram"
	"lessonmode/internal/ui"
	"lessonmode/internal/winutil"
)

//go:embed assets/blocklist.txt
var rawBlocklist string

// Build vaqtida -ldflags bilan kiritilishi mumkin (token manbaga yozilmaydi):
//
//	go build -ldflags "-X main.buildBotToken=XXX -X main.buildAdminID=123"
var (
	buildBotToken string
	buildAdminID  string
)

// unlock kodi holati
type codeState struct {
	mu       sync.Mutex
	code     string
	expires  time.Time
	attempts int
}

func (s *codeState) set(code string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.code = code
	s.expires = time.Now().Add(ttl)
}

func (s *codeState) verify(input string, maxAttempts int) (ok bool, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	input = strings.TrimSpace(input)
	if s.code == "" {
		return false, "Avval \"Ochishni so'rash\" tugmasini bosing."
	}
	if time.Now().After(s.expires) {
		s.code = ""
		return false, "Kod muddati tugadi. Qaytadan so'rang."
	}
	if input != s.code {
		s.attempts++
		return false, "Kod noto'g'ri. Qaytadan urining."
	}
	s.code = ""
	s.attempts = 0
	return true, "Ochildi! Dars rejimi o'chirilmoqda..."
}

func main() {
	args := os.Args[1:]
	switch {
	case hasArg(args, "--uninstall"):
		doUninstall()
		return
	case hasArg(args, "--help"), hasArg(args, "-h"):
		winutil.MessageBox("LessonMode",
			"LessonMode — sinf lockdown dasturi\n\n"+
				"Rejimlar:\n"+
				"  (argumentsiz)  — dars rejimini yoqadi\n"+
				"  --uninstall    — dasturни o'chiradi (admin kerak)\n\n"+
				"Ochish uchun o'quvchi \"Ochishni so'rash\" tugmasini bosadi,\n"+
				"kod ustozning Telegram botiga boradi.", winutil.MBIconInfo)
		return
	default:
		doRun(args)
	}
}

func doRun(args []string) {
	cfg, _ := config.Load()
	applyBuildDefaults(&cfg)

	// Sozlanmagan bo'lsa — qulflamaymiz (tasodifiy "qamalib qolish"ning oldini olish).
	if !cfg.IsConfigured() {
		winutil.MessageBox("LessonMode — sozlash kerak",
			"Bot token va admin ID kiritilmagan.\n\n"+
				"config.json faylini (.exe yonида) to'ldiring:\n"+
				"  botToken: @BotFather dan olingan token\n"+
				"  adminId:  ustozning Telegram ID raqami\n\n"+
				"Yoki build vaqtida -ldflags bilan kiriting.\n"+
				"Hozircha dars rejimi YOQILMADI.", winutil.MBIconWarn)
		return
	}

	// hosts faylini tahrirlash uchun administrator huquqi kerak.
	if !winutil.IsAdmin() {
		if err := winutil.RelaunchElevated(args); err != nil {
			winutil.MessageBox("LessonMode",
				"Administrator huquqi kerak. UAC so'roviga rozilik bering.", winutil.MBIconWarn)
		}
		return
	}

	// Birinchi marta ishga tushsa — o'zini o'rnatadi va logon avtostartини yoqadi.
	if cfg.AutoStart {
		if _, err := installer.Install(); err == nil {
			_ = installer.EnableAutoStart()
		}
	}

	// Bloklanadigan domenlar ro'yxatini tayyorlaymiz va qo'llaymiz.
	domains := blocklist.Build(rawBlocklist, cfg.ExtraBlocked, cfg.ExtraAllowed)
	_ = hosts.Apply(domains)

	// Fon kuzatuvchi: har ReapplySeconds sekundda blokni va ilovalarni tekshiradi.
	stop := make(chan struct{})
	go watchdog(cfg, domains, stop)

	cs := &codeState{}

	handlers := ui.Handlers{
		OnRequestCode: func() string {
			code := telegram.GenerateCode()
			cs.set(code, time.Duration(cfg.CodeTTLSeconds)*time.Second)
			host, _ := os.Hostname()
			user := os.Getenv("USERNAME")
			// Tarmoq so'rovini fon rejimida bajaramiz — oyna qotib qolmasin.
			go func() {
				err := telegram.SendUnlockCode(cfg.BotToken, cfg.AdminID, code, host, user)
				if err != nil {
					winutil.MessageBox("LessonMode", "Kodni yuborishda xato:\n"+err.Error(), winutil.MBIconError)
				}
			}()
			return "Kod ustozingizning Telegram'iga yuborildi.\nUstozdan so'rang va shu yerga kiriting."
		},
		OnVerify: func(code string) (bool, string) {
			ok, msg := cs.verify(code, cfg.MaxAttempts)
			if ok {
				close(stop)
				_ = hosts.Remove()
				_ = installer.DisableAutoStart()
			}
			return ok, msg
		},
	}

	status := cfg.SchoolName + " yoniq — dars rejimi faol"
	hint := "Diqqat jamlangan. Ijtimoiy tarmoq va AI saytlar bloklangan."
	ui.Run(cfg.SchoolName, status, hint, handlers)
}

func watchdog(cfg config.Config, domains []string, stop chan struct{}) {
	ticker := time.NewTicker(time.Duration(cfg.ReapplySeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			if !hosts.IsApplied() {
				_ = hosts.Apply(domains)
			}
			if cfg.CloseDistractingApps {
				procwatch.CloseDistracting(procwatch.DefaultApps)
			}
		}
	}
}

func doUninstall() {
	if !winutil.IsAdmin() {
		if err := winutil.RelaunchElevated([]string{"--uninstall"}); err != nil {
			winutil.MessageBox("LessonMode", "Administrator huquqi kerak.", winutil.MBIconWarn)
		}
		return
	}
	_ = hosts.Remove()
	_ = installer.DisableAutoStart()
	winutil.MessageBox("LessonMode", "LessonMode o'chirildi. Bloklash bekor qilindi.", winutil.MBIconInfo)
}

func applyBuildDefaults(cfg *config.Config) {
	if cfg.BotToken == "" && buildBotToken != "" {
		cfg.BotToken = buildBotToken
	}
	if cfg.AdminID == 0 && buildAdminID != "" {
		if id, err := strconv.ParseInt(strings.TrimSpace(buildAdminID), 10, 64); err == nil {
			cfg.AdminID = id
		}
	}
}

func hasArg(args []string, name string) bool {
	for _, a := range args {
		if strings.EqualFold(a, name) {
			return true
		}
	}
	return false
}
