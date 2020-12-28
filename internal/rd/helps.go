package rd

type Match struct {
	Content  []byte
	Group    map[string]*Match
	Submatch []*Match
}

func (m *Match) Length() int { return len(m.Content) }

func BuildMatch(ms ...interface{}) *Match {
	g := make(map[string]*Match, len(ms)/2)
	c := make([]byte, 0)
	var n string
	for i, x := range ms {
		if i%2 == 0 {
			n = x.(string)
		} else if x.(*Match) != nil {
			if n != "" {
				g[n] = x.(*Match)
			}
			c = append(c, x.(*Match).Content...)
		}
	}

	return &Match{Content: c, Group: g}
}

type Matcher func(cs []byte) (*Match, []byte)

func MatchOne(cs []byte, pred func(c byte) bool) (*Match, []byte) {
	c := cs[0]
	if pred(c) {
		return &Match{Content: cs[0:0]}, cs[:1]
	}

	return nil, nil
}

func SelectLongest(ms map[string]*Match) string {
	var ln string
	var lm *Match

	for n, m := range ms {
		if lm == nil || m.Length() > lm.Length() {
			ln = n
			lm = m
		}
	}

	return ln
}

func MatchOneRune(cs []byte, c rune) (*Match, []byte) {
	return MatchOne(cs, func(b byte) bool { return b == byte(c) })
}

func MatchLongest(cs []byte, ms ...interface{}) (*Match, []byte) {
	msm := make(map[string]*Match, len(ms)/2)
	msr := make(map[string][]byte, len(ms)/2)

	var n string
	for i, x := range ms {
		if i%2 == 0 {
			n = x.(string)
		} else {
			mp := x.(Matcher)
			if m, csr := mp(cs); m != nil {
				msm[n] = m
				msr[n] = csr
			}
		}
	}

	if w := SelectLongest(msm); w != "" {
		return BuildMatch(w, msm[w]), msr[w]
	}

	return nil, nil
}

func MatchManyWithSep(cs []byte, min int, mtch Matcher, sep Matcher) (*Match, []byte) {
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

			totalLen += len(pms[0].Content)
			totalLen += len(pms[1].Content)

			mbs = append(mbs, m)
			ms = append(ms, pms[0], pms[1])

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
		Content:  content,
		Group:    map[string]*Match{},
		Submatch: mbs,
	}

	return m, cs
}

func MatchMany(cs []byte, min int, mtch Matcher) (*Match, []byte) {
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
		Content:  content,
		Group:    map[string]*Match{},
		Submatch: ms,
	}

	return m, cs
}
