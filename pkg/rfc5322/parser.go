package rfc5322

import (
	"github.com/zostay/go-addr/internal/rd"
	"github.com/zostay/go-addr/pkg/rfc5234"
)

const (
	TMailboxList rd.ATag = rd.TLast + iota
	TNameAddr
	TAngleAddr
	TGroup
	TDisplayName
	TAddressList
	TAddrSpec
	TDomainLiteral
	TWords
	TAtom
	TDotAtom
	TFWS
	TCText
	TCContents
	TCContent
	TComment
	TCFWS
	TObsFWS
	TQuotedString
	TObsQP
	TObsPhrase
	TObsAngleAddr
	TObsRoute
	TObsDomainList
	TObsMboxList
	TObsAddrList
	TObsGroupList
	TObsLocalPart
	TObsDomain
)

// address         =   mailbox / group

func MatchAddress(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"mailbox", rd.Matcher(MatchMailbox),
		"group", rd.Matcher(MatchGroup),
	)
}

// mailbox         =   name-addr / addr-spec

func MatchMailbox(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"name-addr", rd.Matcher(MatchNameAddr),
		"addr-spec", rd.Matcher(MatchAddrSpec),
	)
}

// name-addr       =   [display-name] angle-addr

func MatchNameAddr(cs []byte) (*rd.Match, []byte) {
	var (
		dn, aa *rd.Match
	)

	if dn, rcs := MatchDisplayName(cs); dn != nil {
		cs = rcs
	}

	aa, cs = MatchAngleAddr(cs)
	if aa == nil {
		return nil, nil
	}

	return rd.BuildMatch(TNameAddr, "display-name", dn, "angle-addr", aa), cs
}

// angle-addr      =   [CFWS] "<" addr-spec ">" [CFWS] /
//                     obs-angle-addr

func MatchAngleAddr(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"cur", rd.Matcher(MatchCurAngleAddr),
		"obs", rd.Matcher(MatchObsAngleAddr),
	)
}

func MatchCurAngleAddr(cs []byte) (*rd.Match, []byte) {
	var (
		cfws1, la, as, ra, cfws2 *rd.Match
	)

	if cfws1, rcs := MatchCFWS(cs); cfws1 != nil {
		cs = rcs
	}

	la, cs = rd.MatchOneRune(rd.TNone, cs, '<')
	if la == nil {
		return nil, nil
	}

	as, cs = MatchAddrSpec(cs)
	if as == nil {
		return nil, nil
	}

	ra, cs = rd.MatchOneRune(rd.TNone, cs, '>')
	if ra == nil {
		return nil, nil
	}

	if cfws2, rcs := MatchCFWS(cs); cfws2 != nil {
		cs = rcs
	}

	return rd.BuildMatch(TAngleAddr, "", cfws1, "", la, "addr-spec", as, "", ra, "", cfws2), cs
}

// group           =   display-name ":" [group-list] ";" [CFWS]

func MatchGroup(cs []byte) (*rd.Match, []byte) {
	var (
		dn, c, gl, s *rd.Match
	)

	dn, cs = MatchDisplayName(cs)
	if dn == nil {
		return nil, nil
	}

	c, cs = rd.MatchOneRune(rd.TNone, cs, ':')
	if c == nil {
		return nil, nil
	}

	if gl, rcs := MatchGroupList(cs); gl != nil {
		cs = rcs
	}

	s, cs = rd.MatchOneRune(rd.TNone, cs, ';')
	if s == nil {
		return nil, nil
	}

	return rd.BuildMatch(TGroup, "display-name", dn, "", c, "group-list", gl, "", s), cs
}

// display-name    =   phrase

func MatchDisplayName(cs []byte) (*rd.Match, []byte) {
	if p, rcs := MatchPhrase(cs); p != nil {
		return rd.BuildMatch(TDisplayName, "phrase", p), rcs
	}

	return nil, nil
}

// mailbox-list    =   (mailbox *("," mailbox)) / obs-mbox-list

func MatchMailboxList(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"mailbox", rd.Matcher(MatchCurMboxList),
		"obs-mbox-list", rd.Matcher(MatchObsMboxList),
	)
}

func MatchCurMboxList(cs []byte) (*rd.Match, []byte) {
	return rd.MatchManyWithSep(TMailboxList, cs, 1,
		MatchMailbox,
		func(cs []byte) (*rd.Match, []byte) { return rd.MatchOneRune(rd.TNone, cs, ',') },
	)
}

// address-list    =   (address *("," address)) / obs-addr-list

func MatchAddressList(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"address", rd.Matcher(MatchCurAddrList),
		"obs-addr-list", rd.Matcher(MatchObsAddrList),
	)
}

func MatchCurAddrList(cs []byte) (*rd.Match, []byte) {
	return rd.MatchManyWithSep(TAddressList, cs, 1,
		MatchAddress,
		func(cs []byte) (*rd.Match, []byte) { return rd.MatchOneRune(rd.TNone, cs, ',') },
	)
}

// group-list      =   mailbox-list / CFWS / obs-group-list

func MatchGroupList(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"mailbox-list", rd.Matcher(MatchMailboxList),
		"cfws", rd.Matcher(MatchCFWS),
		"obs-group-list", rd.Matcher(MatchObsGroupList),
	)
}

// addr-spec       =   local-part "@" domain

func MatchAddrSpec(cs []byte) (*rd.Match, []byte) {
	var (
		lp, at, d *rd.Match
	)

	lp, cs = MatchLocalPart(cs)
	if lp == nil {
		return nil, nil
	}

	at, cs = rd.MatchOneRune(rd.TNone, cs, '@')
	if at == nil {
		return nil, nil
	}

	d, cs = MatchDomain(cs)
	if d == nil {
		return nil, nil
	}

	return rd.BuildMatch(TAddrSpec, "local-part", lp, "", at, "domain", d), cs
}

// local-part      =   dot-atom / quoted-string / obs-local-part

func MatchLocalPart(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"dot-atom", rd.Matcher(MatchDotAtom),
		"quoted-string", rd.Matcher(MatchQuotedString),
		"obs-local-part", rd.Matcher(MatchObsLocalPart),
	)
}

// domain          =   dot-atom / domain-literal / obs-domain

func MatchDomain(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"dot-atom", rd.Matcher(MatchDotAtom),
		"domain-literal", rd.Matcher(MatchDomainLiteral),
		"obs-domain", rd.Matcher(MatchObsDomain),
	)
}

// domain-literal  =   [CFWS] "[" *([FWS] dtext) [FWS] "]" [CFWS]

func MatchDomainLiteral(cs []byte) (*rd.Match, []byte) {
	var (
		pl, lb, lit, rb *rd.Match
	)
	if pl, rcs := MatchCFWS(cs); pl != nil {
		cs = rcs
	}

	lb, cs = rd.MatchOneRune(rd.TNone, cs, '[')
	if lb == nil {
		return nil, nil
	}

	lit, cs = MatchDomainLiteralLiteral(cs)
	if lit == nil {
		return nil, nil
	}

	rb, cs = rd.MatchOneRune(rd.TNone, cs, ']')
	if rb == nil {
		return nil, nil
	}

	return rd.BuildMatch(
		rd.TLiteral,
		"pre-literal", pl,
		"", lb,
		"literal", lit,
		"", rb,
	), cs
}

func MatchDomainLiteralLiteral(cs []byte) (*rd.Match, []byte) {
	return rd.MatchMany(rd.TNone, cs, 0, MatchDomainLiteralLiteralLiteral)
}

func MatchDomainLiteralLiteralLiteral(cs []byte) (*rd.Match, []byte) {
	var (
		fws, dtext *rd.Match
	)

	if fws, rcs := MatchFWS(cs); fws != nil {
		cs = rcs
	}

	dtext, cs = MatchDText(cs)
	if dtext != nil {
		return nil, nil
	}

	return rd.BuildMatch(rd.TNone, "", fws, "dtext", dtext), cs
}

// dtext           =   %d33-90 /          ; Printable US-ASCII
//                     %d94-126 /         ;  characters not including
//                     obs-dtext          ;  "[", "]", or "\\"

func MatchDText(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"printable-ascii",
		rd.Matcher(func(cs []byte) (*rd.Match, []byte) {
			return rd.MatchOne(rd.TNone, cs, func(c byte) bool { return c >= 0x21 && c <= 0x5a })
		}),
		"other-characters",
		rd.Matcher(func(cs []byte) (*rd.Match, []byte) {
			return rd.MatchOne(rd.TNone, cs, func(c byte) bool { return c >= 0x5e && c <= 0x7e })
		}),
		"obs-dtext",
		rd.Matcher(MatchObsDText),
	)
}

// word            =   atom / quoted-string

func MatchWord(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"atom", rd.Matcher(MatchAtom),
		"quoted-string", rd.Matcher(MatchQuotedString),
	)
}

// phrase          =   1*word / obs-phrase

func MatchPhrase(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"words", rd.Matcher(func(cs []byte) (*rd.Match, []byte) { return rd.MatchMany(TWords, cs, 1, MatchWord) }),
		"obs-phrase", rd.Matcher(MatchObsPhrase),
	)
}

// atext           =   ALPHA / DIGIT /    ; Printable US-ASCII
//                     "!" / "#" /        ;  characters not including
//                     "\$" / "%" /        ;  specials.  Used for atoms.
//                     "&" / "'" /
//                     "*" / "+" /
//                     "-" / "/" /
//                     "=" / "?" /
//                     "^" / "_" /
//                     "`" / "{" /
//                     "|" / "}" /
//                     "~"

func MatchAText(cs []byte) (*rd.Match, []byte) {
	if m, rcs := rfc5234.MatchAlpha(cs); m != nil {
		return m, rcs
	} else if m, rcs := rfc5234.MatchDigit(cs); m != nil {
		return m, rcs
	} else {
		return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool {
			return c == byte('!') || c == byte('#') ||
				c == byte('$') || c == byte('%') ||
				c == byte('&') || c == byte('\'') ||
				c == byte('*') || c == byte('+') ||
				c == byte('-') || c == byte('/') ||
				c == byte('=') || c == byte('?') ||
				c == byte('^') || c == byte('_') ||
				c == byte('`') || c == byte('{') ||
				c == byte('|') || c == byte('}') ||
				c == byte('~')
		})
	}
}

// atom            =   [CFWS] 1*atext [CFWS]

func MatchAtom(cs []byte) (*rd.Match, []byte) {
	var (
		pre, at, post *rd.Match
	)

	if pre, rcs := MatchCFWS(cs); pre != nil {
		cs = rcs
	}

	at, cs = rd.MatchMany(rd.TLiteral, cs, 1, MatchAText)
	if at == nil {
		return nil, nil
	}

	if post, rcs := MatchCFWS(cs); post != nil {
		cs = rcs
	}

	return rd.BuildMatch(TAtom, "pre", pre, "atext", at, "post", post), cs
}

// dot-atom-text   =   1*atext *("." 1*atext)

func MatchDotAtomText(cs []byte) (*rd.Match, []byte) {
	return rd.MatchManyWithSep(rd.TLiteral, cs, 1,
		func(cs []byte) (*rd.Match, []byte) { return rd.MatchMany(rd.TNone, cs, 1, MatchAText) },
		func(cs []byte) (*rd.Match, []byte) { return rd.MatchOneRune(rd.TNone, cs, '.') },
	)
}

// dot-atom        =   [CFWS] dot-atom-text [CFWS]

func MatchDotAtom(cs []byte) (*rd.Match, []byte) {
	var (
		pre, at, post *rd.Match
	)

	if pre, rcs := MatchCFWS(cs); pre != nil {
		cs = rcs
	}

	at, cs = MatchDotAtomText(cs)
	if at == nil {
		return nil, nil
	}

	if post, rcs := MatchCFWS(cs); post != nil {
		cs = rcs
	}

	return rd.BuildMatch(rd.TLiteral, "pre", pre, "dot-atom-text", at, "post", post), cs
}

// FWS             =   ([*WSP CRLF] 1*WSP) /  obs-FWS
//                                        ; Folding white space

func MatchFWS(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"cur-fws", rd.Matcher(MatchCurFWS),
		"obs-fws", rd.Matcher(MatchObsFWS),
	)
}

func MatchCurFWS(cs []byte) (*rd.Match, []byte) {
	var (
		wspcrlf, wsp *rd.Match
	)

	if wspcrlf, rcs := MatchCurFWSPre(cs); wspcrlf != nil {
		cs = rcs
	}

	wsp, cs = rd.MatchMany(rd.TLiteral, cs, 1, rfc5234.MatchWSP)
	if wsp == nil {
		return nil, nil
	}

	return rd.BuildMatch(TFWS, "", wspcrlf, "", wsp), cs
}

func MatchCurFWSPre(cs []byte) (*rd.Match, []byte) {
	var (
		wsp, crlf *rd.Match
	)

	if wsp, rcs := rd.MatchMany(rd.TLiteral, cs, 0, rfc5234.MatchWSP); wsp != nil {
		cs = rcs
	}

	crlf, cs = rfc5234.MatchCRLF(cs)
	if crlf == nil {
		return nil, nil
	}

	return rd.BuildMatch(rd.TLiteral, "", wsp, "", crlf), cs
}

// ctext           =   %d33-39 /          ; Printable US-ASCII
//                     %d42-91 /          ;  characters not including
//                     %d93-126 /         ;  "(", ")", or "\"
//                     obs-ctext

func MatchCText(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"cur-ctext", rd.Matcher(MatchCurCText),
		"obs-ctext", rd.Matcher(MatchObsCText),
	)
}

func MatchCurCText(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(TCText, cs, func(c byte) bool {
		return (c >= 0x21 && c <= 0x27) ||
			(c >= 0x2a && c <= 0x5b) ||
			(c >= 0x5d && c <= 0x7e)
	})
}

// ccontent        =   ctext / quoted-pair / comment

func MatchCContent(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"ctext", rd.Matcher(MatchCText),
		"quoted-pair", rd.Matcher(MatchQuotedPair),
		"comment", rd.Matcher(MatchComment),
	)
}

// comment         =   "(" *([FWS] ccontent) [FWS] ")"

func MatchComment(cs []byte) (*rd.Match, []byte) {
	var (
		lp, cc, rp *rd.Match
	)

	lp, cs = rd.MatchOneRune(rd.TLiteral, cs, '(')
	if lp == nil {
		return nil, nil
	}

	cc, cs = rd.MatchMany(TCContents, cs, 0, func(cs []byte) (*rd.Match, []byte) {
		var (
			fws, cc *rd.Match
		)

		if fws, rcs := MatchFWS(cs); fws != nil {
			cs = rcs
		}

		cc, cs = MatchCContent(cs)
		if cs == nil {
			return nil, nil
		}

		return rd.BuildMatch(rd.TNone, "", fws, "ccontent", cc), cs
	})

	if fws, rcs := MatchFWS(cs); fws != nil {
		cc.Content = append(cc.Content, fws.Content...)
		cs = rcs
	}

	rp, cs = rd.MatchOneRune(rd.TLiteral, cs, ')')
	if rp == nil {
		return nil, nil
	}

	return rd.BuildMatch(TComment, "", lp, "comment-content", cc, "", rp), cs
}

// CFWS            =   (1*([FWS] comment) [FWS]) / FWS

func MatchCFWS(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"cfws-with-comment", rd.Matcher(MatchCFWSWithComment),
		"fws", rd.Matcher(MatchFWS),
	)
}

func MatchCFWSWithComment(cs []byte) (*rd.Match, []byte) {
	var (
		pres, post *rd.Match
	)

	pres, cs = rd.MatchMany(rd.TNone, cs, 1, func(cs []byte) (*rd.Match, []byte) {
		var (
			pre, comment *rd.Match
		)

		if pre, rcs := MatchFWS(cs); pre != nil {
			cs = rcs
		}

		comment, cs = MatchComment(cs)
		if comment == nil {
			return nil, nil
		}

		return rd.BuildMatch(rd.TNone, "pre", pre, "comment", comment), cs
	})
	if pres == nil {
		return nil, nil
	}

	if post, rcs := MatchFWS(cs); post != nil {
		cs = rcs
	}

	return rd.BuildMatch(TCFWS, "pres", pres, "post", post), cs
}

// obs-FWS         =   1*WSP *(CRLF 1*WSP)

func MatchObsFWS(cs []byte) (*rd.Match, []byte) {
	var (
		wsp, crlfs *rd.Match
	)

	wsp, cs = rd.MatchMany(rd.TLiteral, cs, 1, rfc5234.MatchWSP)
	if wsp == nil {
		return nil, nil
	}

	crlfs, cs = rd.MatchMany(rd.TLiteral, cs, 0, func(cs []byte) (*rd.Match, []byte) {
		var (
			crlf, wsps *rd.Match
		)

		crlf, cs = rfc5234.MatchCRLF(cs)
		if crlf == nil {
			return nil, nil
		}

		wsps, cs = rd.MatchMany(rd.TLiteral, cs, 1, rfc5234.MatchWSP)
		if wsps == nil {
			return nil, nil
		}

		return rd.BuildMatch(rd.TLiteral, "crlf", crlf, "wsp", wsps), cs
	})
	if crlfs == nil {
		return nil, nil
	}

	return rd.BuildMatch(TObsFWS, "wsp", wsp, "crlfs", crlfs), cs
}

// qtext           =   %d33 /             ; Printable US-ASCII
//                     %d35-91 /          ;  characters not including
//                     %d93-126 /         ;  "\" or the quote character
//                     obs-qtext

func MatchQText(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"printable-ascii",
		rd.Matcher(func(cs []byte) (*rd.Match, []byte) {
			return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool {
				return c == 0x21 ||
					(c >= 0x23 && c <= 0x5b) ||
					(c >= 0x5d && c <= 0x7e)
			})
		}),
		"obs-qtext",
		rd.Matcher(MatchObsQText),
	)
}

// qcontent        =   qtext / quoted-pair

func MatchQContent(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"qtext", rd.Matcher(MatchQText),
		"quoted-pair", rd.Matcher(MatchQuotedPair),
	)
}

// quoted-string   =   [CFWS]
//                     DQUOTE *([FWS] qcontent) [FWS] DQUOTE
//                     [CFWS]

func MatchQuotedString(cs []byte) (*rd.Match, []byte) {
	var (
		cfws1, ldq, qc, rdq, cfws2 *rd.Match
	)

	if cfws1, rcs := MatchCFWS(cs); cfws1 != nil {
		cs = rcs
	}

	ldq, cs = rfc5234.MatchDQuote(cs)
	if ldq == nil {
		return nil, nil
	}

	qc, cs = rd.MatchMany(rd.TLiteral, cs, 0, func(cs []byte) (*rd.Match, []byte) {
		var (
			fws, qc *rd.Match
		)

		if fws, rcs := MatchFWS(cs); fws != nil {
			cs = rcs
		}

		qc, cs = MatchQContent(cs)
		if qc == nil {
			return nil, nil
		}

		return rd.BuildMatch(rd.TLiteral, "", fws, "qcontent", qc), cs
	})
	if qc == nil {
		return nil, nil
	}

	if fws, rcs := MatchFWS(cs); fws != nil {
		cs = rcs
		qc.Content = append(qc.Content, fws.Content...)
	}

	rdq, cs = rfc5234.MatchDQuote(cs)
	if rdq == nil {
		return nil, nil
	}

	if cfws2, rcs := MatchCFWS(cs); cfws2 != nil {
		cs = rcs
	}

	return rd.BuildMatch(TQuotedString, "", cfws1, "", ldq, "quoted-string", qc, "", rdq, "", cfws2), cs
}

// obs-NO-WS-CTL   =   %d1-8 /            ; US-ASCII control
//                     %d11 /             ;  characters that do not
//                     %d12 /             ;  include the carriage
//                     %d14-31 /          ;  return, line feed, and
//                     %d127              ;  white space characters

func MatchObsNoWSCtl(cs []byte) (*rd.Match, []byte) {
	return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool {
		return (c >= 0x1 && c <= 0x8) ||
			c == 0xb ||
			c == 0xc ||
			(c >= 0xe && c <= 0x1f) ||
			c == 0x7f
	})
}

// obs-ctext       =   obs-NO-WS-CTL

func MatchObsCText(cs []byte) (*rd.Match, []byte) {
	return MatchObsNoWSCtl(cs)
}

// obs-qtext       =   obs-NO-WS-CTL

func MatchObsQText(cs []byte) (*rd.Match, []byte) {
	return MatchObsNoWSCtl(cs)
}

// obs-qp          =   "\\" (%d0 / obs-NO-WS-CTL / LF / CR)

func MatchObsQP(cs []byte) (*rd.Match, []byte) {
	var (
		bs, ch *rd.Match
	)

	bs, cs = rd.MatchOneRune(rd.TLiteral, cs, '\\')
	if bs == nil {
		return nil, nil
	}

	ch, cs = rd.MatchLongest(cs,
		"escaped",
		rd.Matcher(func(cs []byte) (*rd.Match, []byte) {
			return rd.MatchOne(rd.TLiteral, cs, func(c byte) bool {
				return c == 0
			})
		}),
		"obs-no-ws-ctl",
		rd.Matcher(MatchObsNoWSCtl),
		"lf",
		rd.Matcher(rfc5234.MatchLF),
		"cr",
		rd.Matcher(rfc5234.MatchCR),
	)
	if ch == nil {
		return nil, nil
	}

	return rd.BuildMatch(TObsQP, "", bs, "", ch), cs
}

// obs-phrase      =   word *(word / "." / CFWS)

func MatchObsPhrase(cs []byte) (*rd.Match, []byte) {
	var (
		head, tail *rd.Match
	)

	head, cs = MatchWord(cs)
	if head == nil {
		return nil, nil
	}

	tail, cs = rd.MatchMany(rd.TLiteral, cs, 0, func(cs []byte) (*rd.Match, []byte) {
		return rd.MatchLongest(cs,
			"word",
			rd.Matcher(MatchWord),
			"period",
			rd.Matcher(func(cs []byte) (*rd.Match, []byte) { return rd.MatchOneRune(rd.TLiteral, cs, '.') }),
			"cfws",
			rd.Matcher(MatchCFWS),
		)
	})
	if tail == nil {
		return nil, nil
	}

	return rd.BuildMatch(TObsPhrase, "head", head, "tail", tail), cs
}

// quoted-pair     =   ("\\" (VCHAR / WSP)) / obs-qp

func MatchQuotedPair(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"escape",
		rd.Matcher(func(cs []byte) (*rd.Match, []byte) {
			var (
				bs, ch *rd.Match
			)

			bs, cs = rd.MatchOneRune(rd.TLiteral, cs, '\\')
			if bs == nil {
				return nil, nil
			}

			ch, cs = rd.MatchLongest(cs,
				"vchar",
				rfc5234.MatchVChar,
				"wsp",
				rfc5234.MatchWSP,
			)
			if ch == nil {
				return nil, nil
			}

			return rd.BuildMatch(rd.TLiteral, "", bs, "", ch), cs
		}),
		"obs-qp",
		rd.Matcher(MatchObsQP),
	)
}

// obs-angle-addr  =   [CFWS] "<" obs-route addr-spec ">" [CFWS]

func MatchObsAngleAddr(cs []byte) (*rd.Match, []byte) {
	var (
		cfws1, la, rt, as, ra, cfws2 *rd.Match
	)

	if cfws1, rcs := MatchCFWS(cs); cfws1 != nil {
		cs = rcs
	}

	la, cs = rd.MatchOneRune(rd.TLiteral, cs, '<')
	if la == nil {
		return nil, nil
	}

	rt, cs = MatchObsRoute(cs)
	if rt == nil {
		return nil, nil
	}

	as, cs = MatchAddrSpec(cs)
	if as == nil {
		return nil, nil
	}

	ra, cs = rd.MatchOneRune(rd.TLiteral, cs, '>')
	if ra == nil {
		return nil, nil
	}

	if cfws2, rcs := MatchCFWS(cs); cfws2 != nil {
		cs = rcs
	}

	return rd.BuildMatch(TObsAngleAddr, "", cfws1, "", la, "obs-route", rt, "addr-spec", as, "", ra, "", cfws2), cs
}

// obs-route       =   obs-domain-list ":"

func MatchObsRoute(cs []byte) (*rd.Match, []byte) {
	var (
		dl, c *rd.Match
	)

	dl, cs = MatchObsDomainList(cs)
	if dl == nil {
		return nil, nil
	}

	c, cs = rd.MatchOneRune(rd.TLiteral, cs, ':')
	if c == nil {
		return nil, nil
	}

	return rd.BuildMatch(TObsRoute, "obs-domain-list", dl, "", c), cs
}

// obs-domain-list =   *(CFWS / ",") "@" domain
//                     *("," [CFWS] ["@" domain])

func MatchObsDomainList(cs []byte) (*rd.Match, []byte) {
	var (
		bf, at, head, tail *rd.Match
	)

	bf, cs = rd.MatchMany(rd.TLiteral, cs, 0, func(cs []byte) (*rd.Match, []byte) {
		return rd.MatchLongest(cs,
			"cfws",
			rd.Matcher(MatchCFWS),
			"comma",
			rd.Matcher(func(cs []byte) (*rd.Match, []byte) { return rd.MatchOneRune(rd.TLiteral, cs, ',') }),
		)
	})
	if bf == nil {
		return nil, nil
	}

	at, cs = rd.MatchOneRune(rd.TLiteral, cs, '@')
	if at == nil {
		return nil, nil
	}

	head, cs = MatchDomain(cs)
	if head == nil {
		return nil, nil
	}

	tail, cs = rd.MatchMany(rd.TLiteral, cs, 0, func(cs []byte) (*rd.Match, []byte) {
		var (
			c, cfws, at, d *rd.Match
		)

		c, cs = rd.MatchOneRune(rd.TLiteral, cs, ',')
		if c == nil {
			return nil, nil
		}

		if cfws, rcs := MatchCFWS(cs); cfws != nil {
			cs = rcs
		}

		if at, rcs := rd.MatchOneRune(rd.TLiteral, cs, '@'); at != nil {
			cs = rcs
		} else {
			return rd.BuildMatch(rd.TLiteral, "", c, "", cfws), cs
		}

		if d, rcs := MatchDomain(cs); d != nil {
			cs = rcs
		}

		return rd.BuildMatch(rd.TLiteral, "", c, "", cfws, "", at, "", d), cs
	})
	if tail == nil {
		return nil, nil
	}

	return rd.BuildMatch(TObsDomainList, "", bf, "", at, "head", head, "tail", tail), cs
}

// obs-mbox-list   =   *([CFWS] ",") mailbox *("," [mailbox / CFWS])

func MatchObsMboxList(cs []byte) (*rd.Match, []byte) {
	var (
		bf, head, tail *rd.Match
	)

	bf, cs = rd.MatchMany(rd.TLiteral, cs, 0, func(cs []byte) (*rd.Match, []byte) {
		var (
			cfws, c *rd.Match
		)

		if cfws, rcs := MatchCFWS(cs); cfws != nil {
			cs = rcs
		}

		c, cs = rd.MatchOneRune(rd.TLiteral, cs, ',')
		if c == nil {
			return nil, nil
		}

		return rd.BuildMatch(rd.TLiteral, "", cfws, "", c), cs
	})
	if bf == nil {
		return nil, nil
	}

	head, cs = MatchMailbox(cs)
	if head == nil {
		return nil, nil
	}

	tail, cs = rd.MatchMany(rd.TLiteral, cs, 0, func(cs []byte) (*rd.Match, []byte) {
		var (
			c, mb *rd.Match
		)

		c, cs = rd.MatchOneRune(rd.TLiteral, cs, ',')
		if c == nil {
			return nil, nil
		}

		mb, rcs := rd.MatchLongest(cs,
			"mailbox",
			rd.Matcher(MatchMailbox),
			"cfws",
			rd.Matcher(MatchCFWS),
		)
		if mb != nil {
			cs = rcs
		}

		return rd.BuildMatch(rd.TLiteral, "", c, "", mb), cs
	})
	if tail == nil {
		return nil, nil
	}

	return rd.BuildMatch(TObsMboxList, "", bf, "head", head, "tail", tail), cs
}

// obs-addr-list   =   *([CFWS] ",") address *("," [address / CFWS])

func MatchObsAddrList(cs []byte) (*rd.Match, []byte) {
	var (
		bf, head, tail *rd.Match
	)

	bf, cs = rd.MatchMany(rd.TLiteral, cs, 0, func(cs []byte) (*rd.Match, []byte) {
		var (
			cfws, c *rd.Match
		)

		if cfws, rcs := MatchCFWS(cs); cfws != nil {
			cs = rcs
		}

		c, cs = rd.MatchOneRune(rd.TLiteral, cs, ',')
		if c == nil {
			return nil, nil
		}

		return rd.BuildMatch(rd.TLiteral, "", cfws, "", c), cs
	})
	if bf == nil {
		return nil, nil
	}

	head, cs = MatchAddress(cs)
	if head == nil {
		return nil, nil
	}

	tail, cs = rd.MatchMany(rd.TLiteral, cs, 0, func(cs []byte) (*rd.Match, []byte) {
		var (
			c, address *rd.Match
		)

		c, cs = rd.MatchOneRune(rd.TLiteral, cs, ',')
		if c == nil {
			return nil, nil
		}

		address, rcs := rd.MatchLongest(cs,
			"address", rd.Matcher(MatchAddress),
			"cfws", rd.Matcher(MatchCFWS),
		)
		if address != nil {
			cs = rcs
		}

		return rd.BuildMatch(rd.TLiteral, "", c, "address", address), cs
	})
	if tail == nil {
		return nil, nil
	}

	return rd.BuildMatch(TObsAddrList, "", bf, "head", head, "tail", tail), cs
}

// obs-group-list  =   1*([CFWS] ",") [CFWS]

func MatchObsGroupList(cs []byte) (*rd.Match, []byte) {
	var (
		head, tail *rd.Match
	)

	head, cs = rd.MatchMany(rd.TLiteral, cs, 1, func(cs []byte) (*rd.Match, []byte) {
		var (
			cfws, c *rd.Match
		)

		if cfws, rcs := MatchCFWS(cs); cfws != nil {
			cs = rcs
		}

		c, cs = rd.MatchOneRune(rd.TLiteral, cs, ',')
		if c == nil {
			return nil, nil
		}

		return rd.BuildMatch(rd.TLiteral, "", cfws, "", c), cs
	})
	if head == nil {
		return nil, nil
	}

	if tail, rcs := MatchCFWS(cs); tail != nil {
		cs = rcs
	}

	return rd.BuildMatch(TObsGroupList, "head", head, "tail", tail), cs
}

// obs-local-part  =   word *("." word)

func MatchObsLocalPart(cs []byte) (*rd.Match, []byte) {
	return rd.MatchManyWithSep(TObsLocalPart, cs, 1, MatchWord,
		func(cs []byte) (*rd.Match, []byte) {
			return rd.MatchOneRune(rd.TLiteral, cs, '.')
		},
	)
}

// obs-domain      =   atom *("." atom)

func MatchObsDomain(cs []byte) (*rd.Match, []byte) {
	var (
		head, tail *rd.Match
	)

	head, cs = MatchAtom(cs)
	if head == nil {
		return nil, nil
	}

	tail, cs = rd.MatchMany(rd.TLiteral, cs, 0, func(cs []byte) (*rd.Match, []byte) {
		var (
			p, a *rd.Match
		)

		p, cs = rd.MatchOneRune(rd.TLiteral, cs, '.')
		if p == nil {
			return nil, nil
		}

		a, cs = MatchAtom(cs)
		if a == nil {
			return nil, nil
		}

		return rd.BuildMatch(rd.TLiteral, "", p, "atom", a), cs
	})
	if tail == nil {
		return nil, nil
	}

	return rd.BuildMatch(rd.TLiteral, "head", head, "tail", tail), cs
}

// obs-dtext       =   obs-NO-WS-CTL / quoted-pair

func MatchObsDText(cs []byte) (*rd.Match, []byte) {
	return rd.MatchLongest(cs,
		"obs-no-ws-ctl", rd.Matcher(MatchObsNoWSCtl),
		"quoted-pair", rd.Matcher(MatchQuotedPair),
	)
}