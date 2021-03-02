package addr

import (
	"errors"
	"io"
	"mime"
	"strings"

	"github.com/zostay/go-addr/pkg/rd"
	p "github.com/zostay/go-addr/pkg/rfc5322"
)

// This is the list of errors that may be returned while parsing and resolving
// the results of the parse.
var (
	// ErrParseConstruction indicates that the parser was able to match the
	// string, but that object does not produce any output for construction.
	// This error should never occur unless you make a direct call to
	// ApplyActions on a rd.Match tag for which ApplyActions does not produce an
	// object.
	ErrParseConstruction = errors.New("parse construction error")

	// ErrTypeMismatch indicates that the parser was able to match the string,
	// but the type of object produced by that match does not match the pointer
	// provided as the second argument to ApplyActions. The built-in Parse
	// functions will never do this.
	ErrTypeMismatch = errors.New("type mismatch")

	// ErrTypeUnknown indicates that the parser was able to match the string,
	// but then the ApplyActions routine produced an object that it has not been
	// written to handle. If this error occurs, there is a bug in the package.
	ErrTypeUnknown = errors.New("unknown applied type")

	// ErrParse indicates that the parser was unable to match the given input.
	ErrParse = errors.New("unable to parse email address")
)

var (
	// CharsetReader is used to help perform MIME word decoding.
	CharsetReader func(charset string, input io.Reader) (io.Reader, error)
)

// ApplyActions is a low-level function that can be used to transform an
// rd.Match object returned by the Parser into a high-level address object. This
// function is called internally by the Parse functions to produce objects from
// parsed strings.
//
// The second argument must be a pointer to the object that the given rd.Match
// should produce (or it may be set to nil). This method will then evaluate the
// Match and all components of the Match to construct an object. If the
// resulting object is compatible with the pointer, it will be assigned or
// converted to the correct type and then assigned to it. (For example, if the
// parse produced an AddrSpec, but the second argument is a pointer to *Mailbox,
// the AddrSpec will be wrapped in a Mailbox first.)
//
// If the second argument is nil, this method will still construct all the
// objects associated with the match and all components of the match. These can
// be retrieved by walking the tree and checking the m.Made part of each match.
//
// On success, the second argument will be set to the constructed object (unless
// the second argument was nil) and the error will be returned as nil.
//
// On failure, an error is returned and the second argument will not be set. The
// match tree may be fully or partially modified to set Made.
//
// In any case, the tree itself will be unmodified except for assignment to the
// Made field of the match and other components.
func ApplyActions(m *rd.Match, mk interface{}) error {
	if m == nil {
		return ErrParse
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
		return ErrParseConstruction
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
			return ErrTypeMismatch
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
			return ErrTypeMismatch
		}
	case AddressList:
		switch mkv := mk.(type) {
		case *AddressList:
			*mkv = mv
		default:
			return ErrTypeMismatch
		}
	case *Group:
		switch mkv := mk.(type) {
		case *Address:
			*mkv = mv
		case **Group:
			*mkv = mv
		default:
			return ErrTypeMismatch
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
			return ErrTypeMismatch
		}
	case string:
		switch mkv := mk.(type) {
		case *string:
			*mkv = mv
		default:
			return ErrTypeMismatch
		}
	default:
		return ErrTypeUnknown
	}

	return nil
}

func applySubmatchActions(m *rd.Match) {
	if m.Tag == rd.TNone {
		return
	}

	if m.Made != nil {
		return
	}

	for i := range m.Submatch {
		_ = ApplyActions(m.Submatch[i], nil)
	}
}

func applyGroupActions(m *rd.Match) {
	if m.Tag == rd.TNone {
		return
	}

	if m.Made != nil {
		return
	}

	for k := range m.Group {
		_ = ApplyActions(m.Group[k], nil)
	}
}

func decodeMIMEWords(in string) string {
	dec := &mime.WordDecoder{
		CharsetReader: CharsetReader,
	}

	out, err := dec.DecodeHeader(in)
	if err != nil {
		return in
	} else {
		return out
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
			decodeMIMEWords(c),
			strings.TrimSpace(string(m.Content)),
		)
		if err != nil {
			return
		}
	case p.TAngleAddr, p.TObsAngleAddr:
		m.Made = m.Group["addr-spec"].Made
	case p.TGroup:
		var mbl MailboxList
		if ggl := m.Group["group-list"]; ggl != nil {
			if ggl.Made != nil {
				mbl = ggl.Made.(MailboxList)
			} else {
				mbl = MailboxList{}
			}
		} else {
			mbl = MailboxList{}
		}

		m.Made = NewGroupParsed(
			m.Group["display-name"].Made.(string),
			mbl,
			strings.TrimSpace(string(m.Content)),
		)
	case p.TDisplayName:
		m.Made = decodeMIMEWords(strings.TrimSpace(m.Group["phrase"].Made.(string)))
	case p.TMailboxList:
		mailboxes := make(MailboxList, len(m.Submatch))
		for i, mb := range m.Submatch {
			mailboxes[i] = mb.Made.(*Mailbox)
		}
		m.Made = mailboxes
	case p.TObsMboxList:
		gh := m.Group["head"]
		gt := m.Group["tail"]
		mailboxes := make(MailboxList, 1, 1+len(gt.Made.(MailboxList)))
		mailboxes[0] = gh.Made.(*Mailbox)
		mailboxes = append(mailboxes, gt.Made.(MailboxList)...)
		m.Made = mailboxes
	case p.TObsMboxTailList:
		mailboxes := make(MailboxList, 0, len(m.Submatch))
		for _, mbo := range m.Submatch {
			if mb, ok := mbo.Made.(*Mailbox); ok {
				mailboxes = append(mailboxes, mb)
			}
		}
		m.Made = mailboxes
	case p.TObsMboxOptionalList:
		gmb := m.Group["mb"]
		if gmb != nil {
			if mb, ok := gmb.Made.(*Mailbox); ok {
				m.Made = mb
			}
		}
	case p.TAddressList:
		addresses := make(AddressList, len(m.Submatch))
		for i, a := range m.Submatch {
			addresses[i] = a.Made.(Address)
		}
		m.Made = addresses
	case p.TObsAddrList:
		gh := m.Group["head"]
		gt := m.Group["tail"]
		mailboxes := make(AddressList, 1+len(gt.Made.(AddressList)))
		mailboxes[0] = gh.Made.(Address)
		mailboxes = append(mailboxes, gt.Made.(AddressList)...)
		m.Made = mailboxes
	case p.TObsAddrTailList:
		mailboxes := make(AddressList, 0, len(m.Submatch))
		for _, mbo := range m.Submatch {
			if mb, ok := mbo.Made.(Address); ok {
				mailboxes = append(mailboxes, mb)
			}
		}
		m.Made = mailboxes
	case p.TObsAddrOptionalList:
		gmb := m.Group["mb"]
		if gmb != nil {
			if mb, ok := gmb.Made.(Address); ok {
				m.Made = mb
			}
		}
	case p.TObsLocalPart:
		var a strings.Builder
		period := false
		for _, l := range m.Submatch {
			if period {
				a.WriteRune('.')
			}
			a.WriteString(l.Made.(string))
			period = true
		}
		m.Made = a.String()
	case p.TAddrSpec:
		m.Made = NewAddrSpecParsed(
			m.Group["local-part"].Made.(string),
			m.Group["domain"].Made.(string),
			strings.TrimSpace(string(m.Content)),
		)
	case p.TObsDomain:
		var a strings.Builder
		a.WriteString(strings.TrimSpace(m.Group["head"].Made.(string)))
		for _, t := range m.Group["tail"].Submatch {
			a.WriteRune('.')
			a.WriteString(strings.TrimSpace(t.Group["atom"].Made.(string)))
		}
		m.Made = a.String()
	case p.TWords:
		var a strings.Builder
		for _, w := range m.Submatch {
			a.WriteString(w.Made.(string))
		}
		m.Made = a.String()
	case p.TAtom:
		m.Made = string(m.Group["atext"].Content)
	case p.TDotAtom:
		m.Made = m.Group["dot-atom-text"].Made.(string)
	case p.TQuotedString:
		m.Made = string(unquotePairs([]byte(m.Group["quoted-string"].Made.(string))))
	case p.TComment:
		m.Made = string(unquotePairs(m.Group["comment-content"].Content))
	}

	return nil
}

// Keeping this around even though I am not currently using it. I think I will
// need it.
// func unfoldFWS(x []byte) []byte {
// 	bytes.ReplaceAll(x, []byte{'\r', '\n', ' '}, []byte{' '})
// 	bytes.ReplaceAll(x, []byte{'\r', '\n', '\t'}, []byte{'\t'})
// 	return x
// }

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
