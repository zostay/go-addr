package addr

import (
	"strings"

	"github.com/zostay/go-addr/pkg/rfc5322"
)

// Group is the concrete object for holding a named group of email addresses.
type Group struct {
	displayName string
	mailboxList MailboxList
	original    string
}

// DisplayName returns the display name of the group of email addresses.
func (g *Group) DisplayName() string { return g.displayName }

// SetDisplayName updates the display name. It will also clear the original
// string if one is set.
func (g *Group) SetDisplayName(dn string) {
	g.displayName = dn
	g.original = ""
}

// MailboxList returns the slice of mailbox address for this group.
func (g *Group) MailboxList() MailboxList { return g.mailboxList }

// SetMailboxList updates the slice of the mailbox address for this group. It
// will also clear the original string if one is set.
func (g *Group) SetMailboxList(mbs MailboxList) {
	if mbs == nil {
		mbs = MailboxList{}
	}

	g.mailboxList = mbs
	g.original = ""
}

// Address returns the CleanString for the MailboxList.
func (g *Group) Address() string { return g.MailboxList().String() }

// Comment always returns the empty string.
func (g *Group) Comment() string { return "" }

// OriginalString will return the originally parsed string, if that string is
// set. This is useful for roundtripping.
func (g *Group) OriginalString() string {
	return g.original
}

// CleanString returns the canonical version of the group email address string.
func (g *Group) CleanString() string {
	var a strings.Builder
	a.WriteString(g.displayName)
	a.WriteString(": ")
	first := true
	for _, mb := range g.mailboxList {
		if !first {
			a.WriteString(", ")
		}
		a.WriteString(mb.CleanString())
		first = false
	}
	a.WriteString(";")
	return a.String()
}

// String is an alias for CleanString.
func (g *Group) String() string { return g.OriginalString() }

// NewGroupParsed constructs and returns a group email address with an
// associated original string.
func NewGroupParsed(dn string, l MailboxList, o string) *Group {
	return &Group{
		displayName: dn,
		mailboxList: l,
		original:    o,
	}
}

// ParseEmailGroup parses the string as a group email address and returns the
// group found.
//
// This may return a group with no error, a group with an error, or just an
// error depending on the input.
//
// If the entire string is parsed and understood to be a group email address, it
// will return the group and no error.
//
// If the first part of the string is parsed and found to be a group email
// address, it will return what it could parse and also a PartialParseError.
//
// If their is an error parsing the string and no part of a group is found, this
// will return no group object and an error.
func ParseEmailGroup(a string) (*Group, error) {
	m, cs := rfc5322.MatchGroup([]byte(a))

	var group *Group
	err := ApplyActions(m, &group)
	if err != nil {
		return nil, err
	}

	if len(cs) > 0 {
		return group, PartialParseError{string(cs)}
	}

	return group, nil
}
