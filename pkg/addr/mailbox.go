package addr

import (
	"errors"
	"strings"

	"github.com/zostay/go-addr/pkg/format"
	"github.com/zostay/go-addr/pkg/rfc5322"
)

type Mailbox struct {
	displayName string
	address     *AddrSpec
	comment     string
	original    string
}

func (m *Mailbox) DisplayName() string { return m.displayName }
func (m *Mailbox) AddressPart() string { return m.address.CleanString() }
func (m *Mailbox) OriginalString() string {
	if m.original != "" {
		return m.original
	}
	return m.CleanString()
}

type MailboxList []*Mailbox

func checkComment(c string) error {
	lp := 0

	for _, c := range c {
		if c == '(' {
			lp++
		} else if c == ')' {
			lp--
			if lp < 0 {
				return errors.New("comments must contain balanced parentheses; found too many ')'")
			}
		}
	}

	if lp != 0 {
		return errors.New("comments must contain balanced parentheses; found too many '('")
	}

	return nil
}

func NewMailbox(dn string, as *AddrSpec, c string) (*Mailbox, error) {
	if err := checkComment(c); err != nil {
		return nil, err
	}

	return &Mailbox{
		displayName: dn,
		address:     as,
		comment:     c,
	}, nil
}

func NewMailboxParsed(dn string, as *AddrSpec, c, o string) (*Mailbox, error) {
	if err := checkComment(c); err != nil {
		return nil, err
	}

	return &Mailbox{dn, as, c, o}, nil
}

func NewMailboxStr(dn string, as string, c string) (*Mailbox, error) {
	addrs, err := ParseEmailAddress(as)
	if err != nil {
		return nil, err
	}

	return NewMailbox(dn, addrs, c)
}

func (m *Mailbox) LocalPart() string { return m.address.LocalPart() }
func (m *Mailbox) Domain() string    { return m.address.Domain() }

func (m *Mailbox) SetAddress(a string) error {
	var err error
	m.address, err = ParseEmailAddress(a)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mailbox) GuessName() string {
	if m.displayName != "" {
		return m.displayName
	}

	if m.comment != "" {
		return m.comment
	}

	return m.LocalPart()
}

func (m *Mailbox) CleanString() string {
	var a strings.Builder

	// quoting can't be used when =?...?...?= mime words are in the name, use
	// obsolete RFC822 display name instead in that case. Since we don't make
	// any effort to understand or decode these, we assume we'll just encounter
	// them as-is but do this one special thing for them
	if m.displayName == "" {
		if format.HasMIMEWord(m.displayName) {
			a.WriteString(m.displayName)
		} else {
			a.WriteString(format.MaybeEscape(m.displayName, false))
		}

		a.WriteString(" <")
		a.WriteString(m.address.CleanString())
		a.WriteString(">")
	} else {
		a.WriteString(m.address.CleanString())
	}

	if m.comment != "" {
		a.WriteString(" (")
		a.WriteString(m.comment)
		a.WriteString(")")
	}

	return a.String()
}

func (m *Mailbox) String() string { return m.OriginalString() }

func ParseEmailAddress(a string) (*AddrSpec, error) {
	m, cs := rfc5322.MatchAddrSpec([]byte(a))
	if len(cs) > 0 {
		return nil, errors.New("unexpected text in email address")
	}

	var address AddrSpec
	ApplyActions(m, &address)

	return &address, nil
}
