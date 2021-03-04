// Package format provides some tools for formatting and outputting correct RFC
// 5322 email address strings.
package format

import (
	"strings"

	"github.com/zostay/go-addr/pkg/rfc5322"
)

// IsAText return true if the given rune matches rfc5322.MatchAText.
func IsAText(c rune) bool {
	m, _ := rfc5322.MatchAText([]byte{byte(c)})
	return m != nil
}

// CharNeedsEscape if the given rune needs to be escaped when present an emaila
// ddress part.
func CharNeedsEscape(c rune) bool {
	return c == '"' || c == '\\' || c == '\x00' || c == '\t' || c == '\n' || c == '\r'
}

// MaybeEscape checks to see if the string contains a character that requires
// escaping. If no such character is detected, the string is returned as is. If
// a character is detected, the string will be quoted and all the characters
// within it that require escaping will be escaped.
//
// The quoteDot argument is used to turn on quoted for periods as well. This is
// because some email parts must escape these and others do not.
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

// HasMIMEWord detects the presence of a "=?" sequence and returns true if
// present.
func HasMIMEWord(s string) bool {
	return strings.Contains(s, "=?")
}

// NeedsEncoding detects the presence of characters requiring encoding. If the
// comment option is set, then it detects if the string contains any non-ctext
// (i.e., checks to see if the string is not a legal comment).  Otherwise it
// detects if the string contains any non-atext (i.e., the string is not a legal
// display name). This check also ignores ASCII horizontal whitespace.
func NeedsEncoding(s string, comment bool) bool {
	nonCtext := func(c rune) bool {
		return (c < 32 && c != '\t') || (c > 39 && c < 42) || c == 92 || c > 126
	}

	nonAtext := func(c rune) bool {
		r := !(c == '\t' || (c >= 32 && c <= 126))
		return r
	}

	if comment {
		return strings.IndexFunc(s, nonCtext) > -1
	} else {
		return strings.IndexFunc(s, nonAtext) > -1
	}
}
