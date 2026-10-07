# LessonMode

**Sinf uchun diqqatni jamlash (lockdown) dasturi.** Dars paytida o'quvchilar
ijtimoiy tarmoqlar, messenjerlar, video va **AI xizmatlaridan** (ChatGPT, Claude,
Gemini va h.k.) chalg'imasligi uchun ularni bloklaydi. Chiqish faqat **ustozning
Telegram botiga** boradigan 4 xonali kod bilan mumkin.

Bitta `.exe` — Go'da yozilgan, hech qanday runtime yoki tashqi kutubxona kerak emas.

---

## ⚠️ Qonuniy foydalanish (muhim)

Bu dastur **faqat quyidagi shartlarda** ishlatilishi uchun mo'ljallangan:

- Siz boshqarishga **vakolatli** kompyuterlarda (maktab/sinf kompyuterlari);
- O'quvchilar (voyaga yetmaganlarда — ota-onasi) **xabardor** bo'lgan holda.

Dastur **shaffof**: doim ekranda ko'rinib turuvchi oyna bor; administrator uni
istalgan vaqtda `--uninstall` bilan yoki `hosts` faylини tahrirlab olib tashlashi
mumkin. Dastur tizim funksiyalarini (qayta yuklash, Task Manager) **bloklamaydi**
va jarayonни o'ldirishga **qarshilik qilmaydi**. Begona yoki shaxsiy kompyuterларга
yashirin o'rnatish uchun emas.

---

## Nima qiladi

- 🧠 **AI bloklash** — ChatGPT, Claude, Gemini, Copilot, Perplexity, DeepSeek,
  Midjourney va boshqa 120+ AI domeni.
- 📱 **Ijtimoiy tarmoq va messenjer** — Instagram, TikTok, Telegram (web), Facebook,
  Twitter/X, Pinterest, Discord, WhatsApp web va h.k.
- 🎬 **Video/stream va o'yin** — YouTube, Twitch, Netflix, Roblox, Steam va h.k.
- 📋 Jami **1300+ domen** (`assets/blocklist.txt`), `config.json` orqali kengaytiriladi.
- 🪟 **Doim ustda turuvchi** kichik "LessonMode" oynasi (tugma bilan).
- 🔁 Har **5 soniyada** blokni tekshiradi va tiklaydi; chalg'ituvchi ilovalarni
  (Telegram Desktop, o'yinlar...) yopadi — tizimни qotirmasdan.
- 🔑 **Telegram orqali ochish**: o'quvchi tugmani bosadi → 4 xonali kod **ustozning
  botiga** boradi → faqat ustoz aytgan kod bilan ochiladi.

> **Eslatma (halol):** sayt bloklash `hosts` fayli orqali (domen darajasida)
> amalga oshiriladi. Dastur shifrlangan HTTPS trafik ichini ochib o'qimaydi
> (bu uchun MITM kerak bo'lardi — biz buni qilmaymiz). AI aniqlash keng AI-domen
> ro'yxati asosida ishlaydi. Eng kuchli natija uchun buni tarmoq/DNS filtri bilan
> birga ishlating (pastга qarang).

---

## 1-qadam: Telegram botни tayyorlash

1. Telegram'da [@BotFather](https://t.me/BotFather) → `/newbot` → token oling.
   - Agar tokeningiz boshqalarга ko'rinib ketgan bo'lsa — `/revoke` bilan yangilang!
2. O'z Telegram ID raqamingizni biling: [@userinfobot](https://t.me/userinfobot).
3. Botinggiz bilan bir marta `/start` bosing (bot sizга yoza olishi uchun).

## 2-qadam: `.exe` yasash (bitta buyruq)

Go kerak: https://go.dev/dl/

**Windows (PowerShell):**
```powershell
.\build.ps1 -BotToken "SIZNING_TOKEN" -AdminId 123456789
```

**Linux/macOS:**
```bash
BOT_TOKEN=SIZNING_TOKEN ADMIN_ID=123456789 ./build.sh
```

Tayyor fayl: `dist\LessonMode.exe`. Token va ID **exe ichiga** joylashtiriladi —
shu bitta faylni tarqatasiz, `config.json` kerak emas.

> Token/ID bermasdan ham yasashingiz mumkin — u holda `.exe` yonида `config.json`
> turishi kerak (`config.sample.json` dan nusxa oling).

## 3-qadam: Tarqatish va ishga tushirish

- `LessonMode.exe` ni o'quvchilarга bering (Veyon, umumiy papka va h.k.).
- O'quvchi uni ishga tushiradi → **administrator huquqi** so'raladi (UAC) →
  rozilik berilgach:
  - 1300+ sayt bloklanadi,
  - doim ustда turuvchi oyna chiqadi,
  - keyingi safar logonда **avtomatik** ishga tushadi (`autoStart: true`).

> Administrator huquqi `hosts` faylини tahrirlash uchun **zarur**. Veyon yoki domen
> siyosati (GPO) orqali admin bilan tarqatsangiz, UAC avtomatik hal bo'ladi.

## 4-qadam: Ochish (ustoz nazorati)

1. O'quvchi **"Ochishni so'rash"** tugmasini bosadi.
2. 4 xonali kod **ustozning Telegram botiga** keladi (kompyuter nomi bilan).
3. Ustoz kodni aytadi → o'quvchi kiritadi → **"Tasdiqlash"** → dars rejimi o'chadi.

Kod 5 daqiqa amal qiladi (`codeTtlSeconds`). O'quvchi boshqa yo'l bilan ocholmaydi.

---

## Sozlamalar (`config.json`)

| Kalit | Ma'nosi |
|------|---------|
| `botToken` | Telegram bot tokeni |
| `adminId` | Ustozning Telegram ID raqami |
| `schoolName` | Oynada ko'rsatiladigan nom |
| `closeDistractingApps` | Chalg'ituvchi ilovalarni yopish (true/false) |
| `reapplySeconds` | Blokni tekshirish oralig'i (standart 5) |
| `codeTtlSeconds` | Kod muddati, soniya (standart 300) |
| `maxAttempts` | Noto'g'ri urinishlar chegarasi |
| `extraBlocked` | Qo'shimcha bloklanadigan domenlar |
| `extraAllowed` | Hech qachon bloklanmaydigan domenlar (oq ro'yxat) |
| `autoStart` | Logonда avtomatik ishga tushish |

Bloklist ro'yxatini qayta yaratish: `python3 tools/gen_blocklist.py`.

---

## O'chirish (admin uchun)

```powershell
LessonMode.exe --uninstall
```
Bu: `hosts` blokini olib tashlaydi va avtostartني o'chiradi.

**Favqulodda holat** (masalan internet yo'q, kod kelmayapti): administrator
hisobidan `C:\Windows\System32\drivers\etc\hosts` faylини oching va
`# ==== LessonMode BEGIN ...` dan `# ==== LessonMode END ====` gacha bo'lgan
qatorlarни o'chiring, so'ng `ipconfig /flushdns`.

---

## Yanada kuchli: DNS filtri bilan birga

`hosts` bloklash o'quvchi admin bo'lgan kompyuterда aylanib o'tilishi mumkin.
Maksimal himoya uchun **sinf tarmog'ida DNS filtri** qo'shing (masalan NextDNS yoki
router darajasida) — u butun tarmoqда ishlaydi va o'quvchi account'idan buzib
bo'lmaydi. `assets/blocklist.txt` ni o'sha filtrга ham import qilsangiz bo'ladi.

---

## Build (ishlab chiquvchilar uchun)

```bash
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -H=windowsgui" -o dist/LessonMode.exe .
```

GitHub Actions (`.github/workflows/build.yml`) har push'да `.exe` yasaydi
(Artifacts), `v*` teg qo'yilganда esa Release chiqaradi.

## Loyiha tuzilishi

```
main.go                     — kirish nuqtasi, rejimlar, boshqaruv
internal/config/            — config.json
internal/blocklist/         — domen ro'yxatini tayyorlash
internal/hosts/             — hosts faylини boshqarish
internal/procwatch/         — chalg'ituvchi ilovalarни yopish
internal/telegram/          — ustozga kod yuborish
internal/winutil/           — admin tekshiruvi, UAC, xabar oynasi
internal/ui/                — always-on-top Win32 oyna
assets/blocklist.txt        — 1300+ domen (exe ichiga joylashtiriladi)
tools/gen_blocklist.py      — ro'yxat generatori
```

## Litsenziya / mas'uliyat

Faqat vakolatli va xabardor sinf muhitida foydalanish uchun. Dasturdan
foydalanish uchun javobgarlik foydalanuvchining (o'qituvchi/maktab) zimmasida.
