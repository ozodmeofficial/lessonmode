// Package blocklist bloklanadigan domenlar ro'yxatini tayyorlaydi.
package blocklist

import (
	"bufio"
	"strings"
)

// NeverBlock - hech qachon bloklanmaydigan domenlar.
// api.telegram.org MUHIM: ustoz kodни shu orqali oladi, bloklanmasligi shart.
var NeverBlock = map[string]bool{
	"api.telegram.org": true,
	"localhost":        true,
}

// Build embedded ro'yxat (raw) + qo'shimcha/oq ro'yxatdan yakuniy domenlar to'plamini yasaydi.
// Natija tartiblanmagan, lekin takrorlanmaydigan domenlar ro'yxati.
func Build(raw string, extraBlocked, extraAllowed []string) []string {
	allowed := map[string]bool{}
	for k := range NeverBlock {
		allowed[k] = true
	}
	for _, a := range extraAllowed {
		a = normalize(a)
		if a != "" {
			allowed[a] = true
		}
	}

	set := map[string]bool{}
	add := func(d string) {
		d = normalize(d)
		if d == "" || allowed[d] {
			return
		}
		set[d] = true
	}

	sc := bufio.NewScanner(strings.NewReader(raw))
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		add(line)
	}
	for _, e := range extraBlocked {
		add(e)
	}

	out := make([]string, 0, len(set))
	for d := range set {
		out = append(out, d)
	}
	return out
}

// normalize domen nomini tozalaydi: kichik harf, http(s):// va yo'l qismini olib tashlaydi.
func normalize(d string) string {
	d = strings.TrimSpace(strings.ToLower(d))
	d = strings.TrimPrefix(d, "http://")
	d = strings.TrimPrefix(d, "https://")
	if i := strings.IndexAny(d, "/?#"); i >= 0 {
		d = d[:i]
	}
	d = strings.TrimPrefix(d, "www.") // www. variantini alohida qo'shamiz
	d = strings.Trim(d, ".")
	// faqat haqiqiy domenga o'xshaganlarini qoldiramiz
	if d == "" || !strings.Contains(d, ".") {
		return ""
	}
	return d
}
