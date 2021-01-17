package addr

import (
	"errors"
	"strings"

	"github.com/zostay/go-addr/pkg/format"
	"github.com/zostay/go-addr/pkg/rfc5322"
)

var (
	// ErrCommentUnbalancedRight error is returned when constructing a mailbox
	// and the comment given contains too many right parentheses.
	ErrCommentUnbalancedRight = errors.New("comments must contain balanced parentheses; found too many ')'")

	// ErrCommentUnbalancedLeft error is returned when constructing a mailbox
	// and the comment given contains too many left parentheses.
	ErrCommentUnbalancedLeft = errors.New("comments must contain balanced parentheses; found too many '('")
)

// Mailbox is a concrete type for storing a mailbox email address.
type Mailbox struct {
	displayName string
	address     *AddrSpec
	comment     string
	original    string
}

// DisplayName returns the display name of the email address or an empty string.
func (m *Mailbox) DisplayName() string { return m.displayName }

// Address returns the address part of the email address as a string in angle
// brackets. This will return the clean address.
func (m *Mailbox) Address() string { return m.address.String() }

// AddrSpec returns the AddrSpec used to store the email address in detail.
func (m *Mailbox) AddrSpec() *AddrSpec { return m.address }

// Comment returns the accumulated comment for the email address as a string or
// it returns an empty string if there is no comment.
func (m *Mailbox) Comment() string { return m.comment }

// OriginalString either returns the originally parsed string used to create
// this mailbox or an empty string. This is useful for roundtripping, but it
// should not be used for generating new email. See
// https://tools.ietf.org/html/rfc5322#section-4
func (m *Mailbox) OriginalString() string { return m.original }

func checkComment(c string) error {
	lp := 0

	for _, c := range c {
		if c == '(' {
			lp++
		} else if c == ')' {
			lp--
			if lp < 0 {
				return ErrCommentUnbalancedRight
			}
		}
	}

	if lp != 0 {
		return ErrCommentUnbalancedLeft
	}

	return nil
}

// NewMailbox will construct a new mailbox from the given display name,
// AddrSpec, and comment.
//
// This will return ErrCommentUnbalancedRight or ErrCommentUnbalancedLeft if a
// comment is given that contains mismatched parentheses.
//
// On success, returns the constructed mailbox object.
func NewMailbox(
	displayName string,
	addrSpec *AddrSpec,
	comment string,
) (*Mailbox, error) {
	if err := checkComment(comment); err != nil {
		return nil, err
	}

	return &Mailbox{
		displayName: displayName,
		address:     addrSpec,
		comment:     comment,
	}, nil
}

// NewMailboxParsed will construct a new mailbox from the given display name,
// AddrSpec, and comment, and will provide the given original as the
// originally parsed string. This is useful for allowing for the round-tripping
// of an email address.
//
// This will return ErrCommentUnbalancedRight or ErrCommentUnbalancedLeft if a
// comment is given that contains mismatched parantheses.
//
// On success, returns the constructed mailbox object.
func NewMailboxParsed(
	displayName string,
	addrSpec *AddrSpec,
	comment,
	original string,
) (*Mailbox, error) {
	if err := checkComment(comment); err != nil {
		return nil, err
	}

	return &Mailbox{displayName, addrSpec, comment, original}, nil
}

// NewMailboxStr is identical in operation to NewMailbox except that it takes
// the addrSpec argument as a string. This will be parsed into an AddrSpec
// to be added to the object.
//
// This will return the error from ParseEmailAddrSpec if the internal call to
// the parser fails.
//
// If ParseEmailAddrSpec returns a PartialParseError, this method will fail with
// that PartialParseError and fail to construct the object.
//
// This will return ErrCommentUnbalancedRight or ErrCommentUnbalancedLeft if a
// comment is given that contains mismatched parantheses.
//
// On success, returns the constructed mailbox object.
func NewMailboxStr(dn string, as string, c string) (*Mailbox, error) {
	addrs, err := ParseEmailAddrSpec(as)
	if err != nil {
		return nil, err
	}

	return NewMailbox(dn, addrs, c)
}

// LocalPart is syntactic sugar for
//  m.AddrSpec().LocalPart()
func (m *Mailbox) LocalPart() string { return m.address.LocalPart() }

// Domain is syntactic sugar for
//	m.AddrSpec().Domain()
func (m *Mailbox) Domain() string { return m.address.Domain() }

// SetDisplayName will update the display name for the mailbox. This will also
// clear an original string if one is set.
func (m *Mailbox) SetDisplayName(dn string) {
	m.displayName = dn
	m.original = ""
}

// SetComment will update the comment for the mailbox. This will also clear the
// original string if one is set.
func (m *Mailbox) SetComment(c string) {
	m.comment = c
	m.original = ""
}

// SetAddrSpec will update the email address for the mailbox. This will also
// clear the original string if one is set.
func (m *Mailbox) SetAddrSpec(as *AddrSpec) {
	m.address = as
	m.original = ""
}

// SetAddress will change the email address stored. It parses the string using
// ParseAddrSpec and updates the object. This will clear the original string if
// one is set.
//
// This will return the error from ParseEmailAddrSpec if the internal call to
// the parser fails.
//
// If ParseEmailAddrSpec returns a PartialParseError, the address will not be
// set.
//
// On success, returns nil.
func (m *Mailbox) SetAddress(a string) error {
	var err error
	m.address, err = ParseEmailAddrSpec(a)
	if err != nil {
		return err
	}

	m.original = ""

	return nil
}

// GuessName is a helper function aimed at helping you guess what name should be
// associated with the user. The logic works like this:
//
// 1. If a display name is present, it will return the display name.
//
// 2. If a comment is present, it will return the comment.
//
// 3. Fallback to the local part of the email address.
//
// Returns the guessed string.
func (m *Mailbox) GuessName() string {
	if m.displayName != "" {
		return m.displayName
	}

	if m.comment != "" {
		return m.comment
	}

	return m.LocalPart()
}

// CleanString returns a proper RFC 5322 email address. If the originally parsed
// email address was using an obsolete format, this will return the correct
// version according to spec.
func (m *Mailbox) CleanString() string {
	var a strings.Builder

	// quoting can't be used when =?...?...?= mime words are in the name, use
	// obsolete RFC822 display name instead in that case. Since we don't make
	// any effort to understand or decode these, we assume we'll just encounter
	// them as-is but do this one special thing for them
	if m.displayName != "" {
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

// String is an alias for CleanString.
func (m *Mailbox) String() string { return m.CleanString() }

// MailboxList is a slice of mailbox pointers.
type MailboxList []*Mailbox

// OriginalString returns all the email addresses using their original format
// for round-tripping.
//
// Please note, though, that this slice does not store the complete mailbox list
// original, so this will not round-trip the originally parsed content. Any
// superfluous whitespace or comments not directly associated with a mailbox
// address will be lost.
func (ms MailboxList) OriginalString() string {
	var a strings.Builder
	first := true
	for _, m := range ms {
		if !first {
			a.WriteString(", ")
		}
		a.WriteString(m.OriginalString())
		first = false
	}
	return a.String()
}

// CleanString returns the email addresses in the canonical RFC 5322 format
// separated by a comma.
func (ms MailboxList) CleanString() string {
	var a strings.Builder
	first := true
	for _, m := range ms {
		if !first {
			a.WriteString(", ")
		}
		a.WriteString(m.CleanString())
		first = false
	}
	return a.String()
}

// String is an alias for CleanString.
func (ms MailboxList) String() string { return ms.CleanString() }

// ParseEmailMailbox parses the email mailbox string.
//
// It is possible for this method to return a mailbox and an error, or just a
// mailbox with no error, or just an error.
//
// If the parse succeeds and there's no remaining unparsed text, the mailbox
// object will be returned and error will be nil.
//
// If the parse partially succeeds, the mailbox object will be constructed and
// returned and a PartialParseError is returned.
//
// If the parse fails, the error is returned and mailbox is nil.
func ParseEmailMailbox(a string) (*Mailbox, error) {
	m, cs := rfc5322.MatchMailbox([]byte(a))

	var mailbox *Mailbox
	err := ApplyActions(m, &mailbox)
	if err != nil {
		return nil, err
	}

	if len(cs) > 0 {
		return mailbox, PartialParseError{string(cs)}
	}

	return mailbox, nil
}

// ParseEmailMailboxList parses a list of mailbox email addresses. Group
// addresses are not permitted in this parse.
//
// It is possible for this method to return a list of mailboxes and an error, or
// just a list of mailboxes with no error, or just an error.
//
// If the parse success and there's no remaining unparsed text, the mailbox list
// is returned and the error will be nil.
//
// If the parse partially succeeds, the mailbox list will be returned as well as
// a PartialParseError.
//
// If the parse fails, teh error is returned and the returned slice will be nil.
func ParseEmailMailboxList(a string) (MailboxList, error) {
	m, cs := rfc5322.MatchMailboxList([]byte(a))

	var mailboxes MailboxList
	err := ApplyActions(m, &mailboxes)
	if err != nil {
		return nil, err
	}

	if len(cs) > 0 {
		return mailboxes, PartialParseError{string(cs)}
	}

	return mailboxes, nil
}
