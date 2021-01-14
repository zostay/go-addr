package addr

import (
	"bytes"
	"errors"
	"strings"

	"github.com/zostay/go-addr/internal/rd"
	p "github.com/zostay/go-addr/pkg/rfc5322"
)

func ApplyActions(m *rd.Match, mk interface{}) error {
	if m == nil {
		return errors.New("email parse failed")
	}

	applySubmatchActions(m)
	applyGroupActions(m)

	err := applyThisAction(m)
	if err != nil {
		return err
	}

	return passThroughMadeObject(m.Made, mk)
}

func passThroughMadeObject(m, mk interface{}) error {
	if mk == nil {
		return nil
	}

	if m == nil {
		return errors.New("parse construction error")
	}

	switch mv := m.(type) {
	case *Mailbox:
		switch mkv := mk.(type) {
		case *Address:
			*mkv = mv
		case **Mailbox:
			*mkv = mv
		case **AddrSpec:
			*mkv = mv.address
		default:
			return errors.New("type mismatch")
		}
	case MailboxList:
		switch mkv := mk.(type) {
		case *AddressList:
			*mkv = make(AddressList, len(mv))
			for i, v := range mv {
				(*mkv)[i] = v
			}
		case *MailboxList:
			*mkv = mv
		default:
			return errors.New("type mismatch")
		}
	case AddressList:
		switch mkv := mk.(type) {
		case *AddressList:
			*mkv = mv
		default:
			return errors.New("type mismatch")
		}
	case *Group:
		switch mkv := mk.(type) {
		case *Address:
			*mkv = mv
		case **Group:
			*mkv = mv
		default:
			return errors.New("type mismatch")
		}
	case *AddrSpec:
		switch mkv := mk.(type) {
		case *Address:
			*mkv = &Mailbox{
				address:  mv,
				original: mv.original,
			}
		case **Mailbox:
			*mkv = &Mailbox{
				address:  mv,
				original: mv.original,
			}
		case **AddrSpec:
			*mkv = mv
		default:
			return errors.New("type mismatch")
		}
	case string:
		switch mkv := mk.(type) {
		case *string:
			*mkv = mv
		default:
			return errors.New("type mismatch")
		}
	default:
		return errors.New("unknown applied type")
	}

	return nil
}

func applySubmatchActions(m *rd.Match) {
	if m.Tag == rd.TNone {
		return
	}

	for i := range m.Submatch {
		ApplyActions(m.Submatch[i], nil)
	}
}

func applyGroupActions(m *rd.Match) {
	if m.Tag == rd.TNone {
		return
	}

	for k := range m.Group {
		ApplyActions(m.Group[k], nil)
	}
}

func applyThisAction(m *rd.Match) (err error) {
	switch m.Tag {
	case rd.TLiteral:
		m.Made = string(m.Content)
	case p.TNameAddr:
		var dn string
		if m.Group["display-name"] != nil {
			dn = m.Group["display-name"].Made.(string)
		} else {
			dn = ""
		}

		c := accumulateComments(m)
		m.Made, err = NewMailboxParsed(
			dn,
			m.Group["angle-addr"].Made.(*AddrSpec),
			c,
			strings.TrimSpace(string(m.Content)),
		)
		if err != nil {
			return
		}
	case p.TAngleAddr:
		m.Made = m.Group["addr-spec"].Made
	case p.TGroup:
		var mbl MailboxList
		if m.Group["group-list"] != nil {
			mbl = m.Group["group-list"].Made.(MailboxList)
		} else {
			mbl = MailboxList{}
		}

		m.Made = NewGroupParsed(
			m.Group["display-name"].Made.(string),
			mbl,
			strings.TrimSpace(string(m.Content)),
		)
	case p.TDisplayName:
		m.Made = strings.TrimSpace(m.Group["phrase"].Made.(string))
	case p.TMailboxList:
		mailboxes := make(MailboxList, len(m.Submatch))
		for i, mb := range m.Submatch {
			mailboxes[i] = mb.Made.(*Mailbox)
		}
		m.Made = mailboxes
	case p.TAddressList:
		addresses := make(AddressList, len(m.Submatch))
		for i, a := range m.Submatch {
			addresses[i] = a.Made.(Address)
		}
		m.Made = addresses
	case p.TAddrSpec:
		m.Made = NewAddrSpecParsed(
			m.Group["local-part"].Made.(string),
			m.Group["domain"].Made.(string),
			strings.TrimSpace(string(m.Content)),
		)
	case p.TDomainLiteral:
		var a strings.Builder
		a.WriteString(m.Group["pre-literal"].Made.(string))
		a.WriteRune('[')
		a.Write(unfoldFWS(unquotePairs(m.Group["literal"].Content)))
		a.WriteRune(']')
		m.Made = a.String()
	case p.TWords:
		var a strings.Builder
		for _, w := range m.Submatch {
			a.WriteString(w.Made.(string))
		}
		m.Made = a.String()
	case p.TAtom:
		var a strings.Builder
		if m.Group["pre"] != nil {
			a.WriteString(m.Group["pre"].Made.(string))
		}
		a.Write(m.Group["atext"].Content)
		if m.Group["post"] != nil {
			a.WriteString(m.Group["post"].Made.(string))
		}
		m.Made = a.String()
	case p.TDotAtom:
		var a strings.Builder
		if m.Group["pre"] != nil {
			a.WriteString(m.Group["pre"].Made.(string))
		}
		a.WriteString(m.Group["dot-atom-text"].Made.(string))
		if m.Group["pre"] != nil {
			a.WriteString(m.Group["post"].Made.(string))
		}
		m.Made = a.String()
	case p.TQuotedString:
		m.Made = string(unquotePairs([]byte(m.Group["quoted-string"].Made.(string))))
	case p.TComment:
		m.Made = string(unquotePairs(m.Group["comment-content"].Content))
	}

	return nil
}

func unfoldFWS(x []byte) []byte {
	bytes.ReplaceAll(x, []byte{'\r', '\n', ' '}, []byte{' '})
	bytes.ReplaceAll(x, []byte{'\r', '\n', '\t'}, []byte{'\t'})
	return x
}

var (
	quotable = map[byte]struct{}{}
)

func init() {
	qps := []byte{
		' ', '\t', 0x0, 0x1, 0x8,
		0xb, 0xc, 0x7f, '\n', '\r',
	}
	for _, qp := range qps {
		quotable[qp] = struct{}{}
	}
	for qp := byte(0xe); qp <= 0x1f; qp++ {
		quotable[qp] = struct{}{}
	}
}

func unquotePairs(x []byte) []byte {
	output := make([]byte, 0, len(x))
	escaping := false
	for _, c := range x {
		if escaping {
			escaping = false
			if _, ok := quotable[c]; !ok {
				output = append(output, '\\')
			}
			output = append(output, c)
		} else if c == '\\' {
			escaping = true
		} else {
			output = append(output, c)
		}
	}

	return output
}

func accumulateComments(m *rd.Match) string {
	c, _ := accumulateCommentsInner(m)
	return c
}

func accumulateCommentsInner(m *rd.Match) (string, bool) {
	switch m.Tag {
	case p.TCContents:
		return string(m.Content), true
	default:
		cs := make([]string, 0)
		for _, sm := range m.Submatch {
			if c, ok := accumulateCommentsInner(sm); ok {
				cs = append(cs, c)
			}
		}
		if len(cs) > 0 {
			return strings.Join(cs, " "), true
		}
		return "", false
	}
}
