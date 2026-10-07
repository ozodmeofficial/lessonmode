//go:build windows

// Package ui LessonMode'ning doim ustda turuvchi (always-on-top) oynasini chizadi.
// Faqat Win32 API (standart kutubxona) ishlatilgan — tashqi kutubxona yo'q.
package ui

import (
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")

	pRegisterClassExW = user32.NewProc("RegisterClassExW")
	pCreateWindowExW  = user32.NewProc("CreateWindowExW")
	pDefWindowProcW   = user32.NewProc("DefWindowProcW")
	pShowWindow       = user32.NewProc("ShowWindow")
	pUpdateWindow     = user32.NewProc("UpdateWindow")
	pGetMessageW      = user32.NewProc("GetMessageW")
	pTranslateMessage = user32.NewProc("TranslateMessage")
	pDispatchMessageW = user32.NewProc("DispatchMessageW")
	pPostQuitMessage  = user32.NewProc("PostQuitMessage")
	pDestroyWindow    = user32.NewProc("DestroyWindow")
	pSetWindowPos     = user32.NewProc("SetWindowPos")
	pSendMessageW     = user32.NewProc("SendMessageW")
	pGetWindowTextW   = user32.NewProc("GetWindowTextW")
	pSetWindowTextW   = user32.NewProc("SetWindowTextW")
	pLoadCursorW      = user32.NewProc("LoadCursorW")
	pGetSystemMetrics = user32.NewProc("GetSystemMetrics")

	pGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
	pGetStockObject   = gdi32.NewProc("GetStockObject")
)

const (
	wsOverlapped  = 0x00000000
	wsCaption     = 0x00C00000
	wsVisible     = 0x10000000
	wsChild       = 0x40000000
	wsBorder      = 0x00800000
	wsTabstop     = 0x00010000
	wsGroup       = 0x00020000
	esNumber      = 0x2000
	esCenter      = 0x0001
	esAutoHScroll = 0x0080

	wsExTopmost    = 0x00000008
	wsExToolWindow = 0x00000080
	wsExClientEdge = 0x00000200

	swShow = 5

	wmDestroy = 0x0002
	wmClose   = 0x0010
	wmCommand = 0x0111
	wmSetFont = 0x0030

	swpNoSize     = 0x0001
	swpNoMove     = 0x0002
	swpShowWindow = 0x0040

	defaultGUIFont = 17
	smCxScreen     = 0
	smCyScreen     = 1

	idRequest = 1001
	idVerify  = 1003
)

// Handlers UI tugmalari bosilganda chaqiriladigan funksiyalar.
type Handlers struct {
	// Ochish kodini so'rash (Telegram ga yuborish). Qaytaradi: ko'rsatiladigan xabar.
	OnRequestCode func() string
	// Kodni tekshirish. ok=true bo'lsa dastur ochiladi (oyna yopiladi).
	OnVerify func(code string) (ok bool, message string)
}

type uiState struct {
	hEdit    uintptr
	hMsg     uintptr
	handlers Handlers
	unlocked bool
}

var state *uiState
var hwndTopmost = ^uintptr(0) // HWND_TOPMOST = -1

type wndClassExW struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     uintptr
	hIcon         uintptr
	hCursor       uintptr
	hbrBackground uintptr
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       uintptr
}

type msg struct {
	hwnd    uintptr
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      struct{ x, y int32 }
}

func u16(s string) *uint16 { p, _ := syscall.UTF16PtrFromString(s); return p }

func setFont(h uintptr) {
	font, _, _ := pGetStockObject.Call(defaultGUIFont)
	pSendMessageW.Call(h, wmSetFont, font, 1)
}

func readEdit() string {
	buf := make([]uint16, 16)
	n, _, _ := pGetWindowTextW.Call(state.hEdit, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return syscall.UTF16ToString(buf[:n])
}

func setMsg(text string) {
	pSetWindowTextW.Call(state.hMsg, uintptr(unsafe.Pointer(u16(text))))
}

func wndProc(hwnd, message, wParam, lParam uintptr) uintptr {
	switch message {
	case wmCommand:
		id := wParam & 0xFFFF
		switch id {
		case idRequest:
			if state.handlers.OnRequestCode != nil {
				setMsg(state.handlers.OnRequestCode())
			}
			return 0
		case idVerify:
			code := readEdit()
			if state.handlers.OnVerify != nil {
				ok, m := state.handlers.OnVerify(code)
				setMsg(m)
				if ok {
					state.unlocked = true
					pDestroyWindow.Call(hwnd)
				}
			}
			return 0
		}
	case wmClose:
		// Oddiy yopishga ruxsat bermaymiz — chiqish uchun ustoz kodi kerak.
		// (Jarayonni Task Manager bilan to'xtatish mumkin — biz bunga qarshilik qilmaymiz.)
		setMsg("Chiqish uchun ustoz bergan kodni kiriting.")
		return 0
	case wmDestroy:
		pPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := pDefWindowProcW.Call(hwnd, message, wParam, lParam)
	return r
}

// Run always-on-top oynani ko'rsatadi va xabar siklini boshlaydi.
// Oyna yopilgach (kod to'g'ri kiritilганda) qaytaradi. unlocked=true bo'lsa — ochildi.
func Run(title, statusLine, hintLine string, handlers Handlers) bool {
	state = &uiState{handlers: handlers}

	hInst, _, _ := pGetModuleHandleW.Call(0)
	cursor, _, _ := pLoadCursorW.Call(0, 32512) // IDC_ARROW
	className := u16("LessonModeWindow")

	wc := wndClassExW{
		style:         0x0002 | 0x0001, // CS_HREDRAW | CS_VREDRAW
		lpfnWndProc:   syscall.NewCallback(wndProc),
		hInstance:     hInst,
		hCursor:       cursor,
		hbrBackground: 6, // COLOR_WINDOW+1
		lpszClassName: className,
	}
	wc.cbSize = uint32(unsafe.Sizeof(wc))
	pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	width, height := 380, 210
	scrW, _, _ := pGetSystemMetrics.Call(smCxScreen)
	// Ekranning yuqori-o'ng burchagiga joylashtiramiz
	x := int(scrW) - width - 20
	y := 20

	exStyle := uintptr(wsExTopmost | wsExToolWindow)
	style := uintptr(wsCaption | wsBorder | wsVisible | wsOverlapped)

	hwnd, _, _ := pCreateWindowExW.Call(
		exStyle,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(u16(title))),
		style,
		uintptr(x), uintptr(y), uintptr(width), uintptr(height),
		0, 0, hInst, 0,
	)
	if hwnd == 0 {
		return false
	}

	// Status matni
	createStatic(hwnd, hInst, statusLine, 15, 12, 345, 22)
	createStatic(hwnd, hInst, hintLine, 15, 38, 345, 20)

	// "Ochishni so'rash" tugmasi
	createButton(hwnd, hInst, "Ochishni so'rash", idRequest, 15, 70, 345, 30)

	// "Kod:" + input + "Tasdiqlash"
	createStatic(hwnd, hInst, "Ustoz bergan kod:", 15, 115, 110, 22)
	state.hEdit = createEdit(hwnd, hInst, 130, 112, 90, 26)
	createButton(hwnd, hInst, "Tasdiqlash", idVerify, 230, 110, 130, 30)

	// Xabar qatori
	state.hMsg = createStatic(hwnd, hInst, "", 15, 150, 345, 36)

	// Doim ustda tursin
	pSetWindowPos.Call(hwnd, hwndTopmost, 0, 0, 0, 0, swpNoMove|swpNoSize|swpShowWindow)
	pShowWindow.Call(hwnd, swShow)
	pUpdateWindow.Call(hwnd)

	// Xabar sikli
	var m msg
	for {
		r, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
	return state.unlocked
}

func createStatic(parent, hInst uintptr, text string, x, y, w, h int) uintptr {
	h2, _, _ := pCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(u16("STATIC"))),
		uintptr(unsafe.Pointer(u16(text))),
		uintptr(wsChild|wsVisible),
		uintptr(x), uintptr(y), uintptr(w), uintptr(h),
		parent, 0, hInst, 0,
	)
	setFont(h2)
	return h2
}

func createButton(parent, hInst uintptr, text string, id, x, y, w, h int) uintptr {
	h2, _, _ := pCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(u16("BUTTON"))),
		uintptr(unsafe.Pointer(u16(text))),
		uintptr(wsChild|wsVisible|wsTabstop),
		uintptr(x), uintptr(y), uintptr(w), uintptr(h),
		parent, uintptr(id), hInst, 0,
	)
	setFont(h2)
	return h2
}

func createEdit(parent, hInst uintptr, x, y, w, h int) uintptr {
	h2, _, _ := pCreateWindowExW.Call(
		uintptr(wsExClientEdge),
		uintptr(unsafe.Pointer(u16("EDIT"))),
		uintptr(unsafe.Pointer(u16(""))),
		uintptr(wsChild|wsVisible|wsTabstop|esNumber|esCenter|esAutoHScroll),
		uintptr(x), uintptr(y), uintptr(w), uintptr(h),
		parent, 0, hInst, 0,
	)
	setFont(h2)
	return h2
}
