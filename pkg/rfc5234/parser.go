package rfc5234

import (
	"github.com/zostay/go-addr/pkg/rd"
)

// ALPHA          =  %x41-5A / %x61-7A   ; A-Z / a-z

func MatchAlpha(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool {
		return (c >= 0x41 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a)
	})
}

// DIGIT          =  %x30-39
//                        ; 0-9

func MatchDigit(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c >= 0x30 && c <= 0x39 })
}

// CR             =  %x0D
//                        ; carriage return

func MatchCR(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c == 0xd })
}

// LF             =  %x0A
//                        ; linefeed

func MatchLF(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c == 0xa })
}

// CRLF           =  CR LF
//                        ; Internet standard newline

func MatchCRLF(cs []byte) (*rd.Match, []byte) {
	var (
		cr, lf *rd.Match
	)

	cr, cs = MatchCR(cs)
	if cr == nil {
		return nil, nil
	}

	lf, cs = MatchLF(cs)
	if lf == nil {
		return nil, nil
	}

	return rd.BuildMatch(rd.TLiteral, "", cr, "", lf), cs
}

// DQUOTE         =  %x22
// 					; " (Double Quote)

func MatchDQuote(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c == 0x22 })
}

// HTAB           =  %x09
//                   ; horizontal tab

func MatchHTab(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c == 0x9 })
}

// SP             =  %x20

func MatchSP(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c == 0x20 })
}

// WSP            =  SP / HTAB
//                        ; white space

func MatchWSP(cs []byte) (*rd.Match, []byte) {
	if m, rcs := MatchSP(cs); m != nil {
		return m, rcs
	} else if m, rcs := MatchHTab(cs); m != nil {
		return m, rcs
	} else {
		return nil, nil
	}
}

// VCHAR          =  %x21-7E
//                        ; visible (printing) characters

func MatchVChar(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c >= 0x21 && c <= 0x7e })
}
