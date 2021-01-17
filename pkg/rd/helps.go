// Package rd provides recursive descent parser helpers. This is not really
// intended of use outside of this library, but some objects here are exposed
// and could be useful in some applications. This should probably be a separate
// package library.
package rd

import (
	"log"
	"reflect"
	"runtime"
)

// Match is the object used to represent some segment of a parsed string.
type Match struct {
	Tag      ATag              // an identifier describing what the match represents
	Content  []byte            // the full content of the match
	Group    map[string]*Match // identifies named submatches
	Submatch []*Match          // identifies a list of submatches
	Made     interface{}       // a place to put high-level objects generated from this match
}

// ATag is the type used to tag matches by type.
type ATag int

// A few standard tags for matches.
const (
	TNone ATag = iota
	TLiteral
	TLast
)

const trace = false

func traceMatch(fmt string, args ...interface{}) {
	if trace {
		log.Printf(fmt, args...)
	}
}

// Length returns the number of bytes matched for this match.
func (m *Match) Length() int {
	if m != nil {
		return len(m.Content)
	} else {
		return 0
	}
}

// BuildMatch is a short hand for building a match with named submatches.
func BuildMatch(t ATag, ms ...interface{}) (m *Match) {
	g := make(map[string]*Match, len(ms)/2)
	s := make([]*Match, 0, len(ms)/2)
	c := make([]byte, 0)
	var n string
	for i, x := range ms {
		if i%2 == 0 {
			n = x.(string)
		} else if x.(*Match) != nil {
			if n != "" {
				g[n] = x.(*Match)
			}
			s = append(s, x.(*Match))
			c = append(c, x.(*Match).Content...)
		}
	}

	m = &Match{Tag: t, Content: c, Group: g, Submatch: s}
	//traceMatch("BuildMatch(%+v)", m)

	return
}

// Matcher is the type for matching functions. These accept a list of bytes to
// start matching from and return a pointer to a Match and a list of remaining
// unmatched bytes.
//
// If the match is successful, then a pointer to a Match is returned and the
// remaining input is also returned. It is possible for a match to match zero
// bytes.
//
// If the match fails, then the Match should be returned as nil. Usually the
// remaining input is also returned as nil in that case.
type Matcher func(cs []byte) (*Match, []byte)

// MatchOne matches exactly one byte if the next byte in the input matches the
// given predicate.
func MatchOne(t ATag, cs []byte, pred func(c byte) bool) (*Match, []byte) {
	if len(cs) == 0 {
		//traceMatch("!MatchOne(%d, %v, %s)", t, string(cs), runtime.FuncForPC(reflect.ValueOf(pred).Pointer()).Name())
		traceMatch("TRY MatchOne(%d, empty, %s)", t, runtime.FuncForPC(reflect.ValueOf(pred).Pointer()).Name())
		return nil, nil
	}

	traceMatch("TRY MatchOne(%d, %v, %s)", t, string(cs), runtime.FuncForPC(reflect.ValueOf(pred).Pointer()).Name())

	c := cs[0]
	if pred(c) {
		m := Match{Tag: t, Content: cs[0:1]}
		traceMatch("GOT MatchOne(%d, %v, %s) = %v", t, string(cs), runtime.FuncForPC(reflect.ValueOf(pred).Pointer()).Name(), m)
		return &m, cs[1:]
	}

	//traceMatch("!MatchOne(%d, %v, %s)", t, string(cs), runtime.FuncForPC(reflect.ValueOf(pred).Pointer()).Name())
	return nil, nil
}

// MatchOneRune matches the next byte if it exactly matches the given rune.
func MatchOneRune(t ATag, cs []byte, c rune) (*Match, []byte) {
	traceMatch("TRY MatchOneRune(%d, %v, %c)", t, string(cs), c)
	return MatchOne(t, cs, func(b byte) bool { return b == byte(c) })
}

func selectLongest(ms []*Match) int {
	var ln int
	var lm *Match

	for n, m := range ms {
		if lm == nil || m.Length() > lm.Length() {
			ln = n
			lm = m
		}
	}

	return ln
}

// MatchLongest tries all the given matchers against the current input. It then
// returns whichever of these matches works to match the most input.
func MatchLongest(cs []byte, ms ...Matcher) (*Match, []byte) {
	msm := make([]*Match, len(ms))
	msr := make([][]byte, len(ms))

	for i, mp := range ms {
		if m, csr := mp(cs); m != nil {
			msm[i] = m
			msr[i] = csr
		}
	}

	if w := selectLongest(msm); w != -1 {
		traceMatch("GOT MatchLongest(%v) = (%d, %v)", string(cs), w, msm[w])
		return msm[w], msr[w]
	}

	return nil, nil
}

// MatchManyWithSep matches the given matcher against the input provided that
// the separator matcher matches in between. It returns a match containing those
// matches. If fewer than min matches are present, the match returns no match.
func MatchManyWithSep(t ATag, cs []byte, min int, mtch Matcher, sep Matcher) (*Match, []byte) {
	mbs := make([]*Match, 0)
	ms := make([]*Match, 0)
	totalLen := 0

	for {
		var pms [2]*Match
		if len(ms) > 0 {
			if m, rcs := sep(cs); m != nil {
				pms[0] = m
				cs = rcs
			} else {
				break
			}
		}
		if m, rcs := mtch(cs); m != nil {
			pms[1] = m
			cs = rcs

			if len(ms) > 0 {
				totalLen += len(pms[0].Content)
			}
			totalLen += len(pms[1].Content)

			mbs = append(mbs, m)
			if len(ms) > 0 {
				ms = append(ms, pms[0], pms[1])
			} else {
				ms = append(ms, pms[1])
			}

			continue
		}

		break
	}

	if len(mbs) < min {
		return nil, nil
	}

	content := make([]byte, 0, totalLen)
	for _, m := range ms {
		content = append(content, m.Content...)
	}

	m := &Match{
		Tag:      t,
		Content:  content,
		Group:    map[string]*Match{},
		Submatch: mbs,
	}

	traceMatch("GOT MatchManyWithSep(%d, %v, %d, %s, %s) = %v",
		t, string(cs), min,
		runtime.FuncForPC(reflect.ValueOf(mtch).Pointer()).Name(),
		runtime.FuncForPC(reflect.ValueOf(sep).Pointer()).Name(),
		m,
	)
	return m, cs
}

// MatchMany matches the given matcher as many times as possible one after
// another on the input. If the number of matches is fewer than min, it returns
// a failure.
func MatchMany(t ATag, cs []byte, min int, mtch Matcher) (*Match, []byte) {
	content := make([]byte, 0)
	ms := make([]*Match, 0)

	for {
		if m, rcs := mtch(cs); m != nil {
			cs = rcs

			ms = append(ms, m)
			content = append(content, m.Content...)

			continue
		}

		break
	}

	if len(ms) < min {
		return nil, nil
	}

	m := &Match{
		Tag:      t,
		Content:  content,
		Group:    map[string]*Match{},
		Submatch: ms,
	}

	traceMatch("GOT MatchMany(%d, %v, %d, %s) = %v",
		t, string(cs), min,
		runtime.FuncForPC(reflect.ValueOf(mtch).Pointer()).Name(),
		m,
	)

	return m, cs
}
