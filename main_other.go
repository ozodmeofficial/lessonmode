//go:build !windows

package main

import "fmt"

// LessonMode faqat Windows uchun mo'ljallangan.
// Bu fayl boshqa platformalarda kompilyatsiya qilish uchun (masalan tekshiruv).
func main() {
	fmt.Println("LessonMode faqat Windows'da ishlaydi. Build: GOOS=windows go build")
}
