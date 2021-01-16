// Package rfc5234 implements an RFC 5234 parser which provides basic
// productions used by the RFC 5322 parser.
package rfc5234

import (
	"github.com/zostay/go-addr/pkg/rd"
)

// MatchAlpha matches a single ASCII alphabetical character.
//  // ALPHA          =  %x41-5A / %x61-7A   ; A-Z / a-z
func MatchAlpha(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool {
		return (c >= 0x41 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a)
	})
}

// MatchDigit matches a single ASCII digit.
//  // DIGIT          =  %x30-39
//  //                        ; 0-9
func MatchDigit(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c >= 0x30 && c <= 0x39 })
}

// MatchCR matches a single carriage return.
//  // CR             =  %x0D
//  //                        ; carriage return
func MatchCR(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c == 0xd })
}

// MatchLF matches a single newline character.
//  // LF             =  %x0A
//  //                        ; linefeed
func MatchLF(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c == 0xa })
}

// MatchCRLF matches a network new line: a carriage return followed by a line
// feed.
//  // CRLF           =  CR LF
//  //                        ; Internet standard newline
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

// MatchDQuote matches a single double quote.
//  // DQUOTE         =  %x22
//  // 					; " (Double Quote)
func MatchDQuote(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c == 0x22 })
}

// MatchHTab matches a single tab.
//  // HTAB           =  %x09
//  //                   ; horizontal tab
func MatchHTab(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c == 0x9 })
}

// MatchSP matches a single space.
//  // SP             =  %x20
func MatchSP(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c == 0x20 })
}

// MatchWSP matches either a single space or tab.
//  // WSP            =  SP / HTAB
//  //                        ; white space
func MatchWSP(cs []byte) (*rd.Match, []byte) {
	if m, rcs := MatchSP(cs); m != nil {
		return m, rcs
	} else if m, rcs := MatchHTab(cs); m != nil {
		return m, rcs
	} else {
		return nil, nil
	}
}

// MatchVChar matches a single visible ASCII character.
//  // VCHAR          =  %x21-7E
//  //                        ; visible (printing) characters
func MatchVChar(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool { return c >= 0x21 && c <= 0x7e })
}
