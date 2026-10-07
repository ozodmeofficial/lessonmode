<#
  LessonMode — bitta buyruq bilan .exe yasash (Windows, PowerShell)

  Foydalanish:
    1) Go o'rnatilgan bo'lsin: https://go.dev/dl/  (go version bilan tekshiring)
    2) PowerShell'ni oching, loyiha papkasiga kiring va ishga tushiring:

         .\build.ps1 -BotToken "SIZNING_TOKEN" -AdminId 123456789

    3) Tayyor fayl:  dist\LessonMode.exe
       Shu bitta faylni o'quvchilarга tarqatasiz (config.json KERAK EMAS,
       token va ID exe ichiga joylashtiriladi).

  Token/ID bermasangiz, exe config.json dan o'qiydi (exe yonида turishi kerak).
#>
param(
    [string]$BotToken = "",
    [string]$AdminId = "",
    [string]$Output = "dist\LessonMode.exe"
)

$ErrorActionPreference = "Stop"

# Go bor-yo'qligini tekshiramiz
try { go version | Out-Null } catch {
    Write-Host "XATO: Go topilmadi. https://go.dev/dl/ dan o'rnating." -ForegroundColor Red
    exit 1
}

$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

$ldflags = "-s -w -H=windowsgui"
if ($BotToken -ne "") { $ldflags += " -X main.buildBotToken=$BotToken" }
if ($AdminId -ne "")  { $ldflags += " -X main.buildAdminID=$AdminId" }

New-Item -ItemType Directory -Force -Path (Split-Path $Output) | Out-Null

Write-Host "LessonMode build qilinmoqda..." -ForegroundColor Cyan
go build -ldflags $ldflags -o $Output .

if (Test-Path $Output) {
    $size = [math]::Round((Get-Item $Output).Length / 1MB, 1)
    Write-Host ""
    Write-Host "TAYYOR: $Output ($size MB)" -ForegroundColor Green
    if ($BotToken -eq "" -or $AdminId -eq "") {
        Write-Host "Eslatma: token/ID berilmadi -> exe yonида config.json turishi kerak." -ForegroundColor Yellow
    } else {
        Write-Host "Token va admin ID exe ichiga joylashtirildi. Bitta faylni tarqatish mumkin." -ForegroundColor Green
    }
} else {
    Write-Host "Build muvaffaqiyatsiz." -ForegroundColor Red
    exit 1
}
