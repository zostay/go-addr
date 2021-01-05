package format

import (
	"strings"

	"github.com/zostay/go-addr/pkg/rfc5322"
)

func IsAText(c rune) bool {
	m, _ := rfc5322.MatchAText([]byte{byte(c)})
	return m != nil
}

func CharNeedsEscape(c rune) bool {
	return c == '"' || c == '\\' || c == '\x00' || c == '\t' || c == '\n' || c == '\r'
}

func MaybeEscape(s string, quoteDot bool) string {
	if s == "" {
		return ""
	}

	// leading or trailing dot is always quoted
	if strings.HasPrefix(s, ".") || strings.HasSuffix(s, ".") {
		quoteDot = true
	}

	quote := false

	// is quoting needed otherwise?
	var a strings.Builder
	a.WriteRune('"')
	for _, c := range s {
		if !IsAText(c) && (quoteDot || c != '.') {
			quote = true
		}

		if CharNeedsEscape(c) {
			quote = true

			a.WriteRune('\\')
			a.WriteRune(c)
		} else {
			a.WriteRune(c)
		}
	}
	a.WriteRune('"')

	if quote {
		return a.String()
	}

	s = a.String()
	return s[1 : len(s)-1]
}

func HasMIMEWord(s string) bool {
	return strings.Contains(s, "=?")
}
