// Package addr is the primary entry point for working with email addresses.
//
// Roundtripping
//
// Generally speaking, roundtripping is not desireable when it comes to email
// addresses. If you are parsing an old email to harvest the email addresses
// from it, you want to make sure any new email being created to those email
// addresses use the current format rather than any obsolete format that may
// have been in use. Use the String() or CleanString() methods to get the
// canonical string for use with new messages and documents.
//
// However, there are caes where being able to understand the details of the
// address and still be able to output the original is still what is wanted.
// This library is capable of helping you to a limited extent in preserving the
// original strings. Use the OriginalString() methods to retrieve the original
// string for roundtripping.
//
// There is a caveat though...
//
// The Mailbox, AddrSpec, and Group objects all store the originally parsed text
// when they are created by the parser. However, the lists of email addresses
// (either AddressList or MailboxList) are not quite able to totally preserve
// the original as these are just slices. Any extra comments or whitespace
// not associated with a mailbox or group email address in the originally parsed
// text will be lost by those data structures.
package addr

import (
	"strings"

	"github.com/zostay/go-addr/pkg/rfc5322"
)

// Address represents a generic email address. This could be either a mailbox
// address or a group address.
type Address interface {
	DisplayName() string    // name associated with the group
	Address() string        // mailbox address(es) as a string
	OriginalString() string // the originally parsed string for round-tripping
	CleanString() string    // a clean string if you want the canonical email string
	Comment() string        // the comment associated with this address
}

// AddressList is a slice of Address objects.
type AddressList []Address

// OriginalString is not entirely a reproduction of the original. An AddressList
// is just a slice of Address, which means it does not preserve the original in
// any way. The addresses will be identical to the parsed original, though. It
// will join these with a comma followed by a space.
//
// Any address in the list that was not parsed or has no original string will
// actually be the canonical string instead.
func (as AddressList) OriginalString() string {
	var a strings.Builder
	first := true
	for _, addr := range as {
		if !first {
			a.WriteString(", ")
		}
		a.WriteString(addr.OriginalString())
		first = false
	}
	return a.String()
}

// CleanString returns the canonical version of the addresses in the list joined
// together with a comma and a space.
func (as AddressList) CleanString() string {
	var a strings.Builder
	first := true
	for _, addr := range as {
		if !first {
			a.WriteString(", ")
		}
		a.WriteString(addr.CleanString())
		first = false
	}
	return a.String()
}

// String is an alias for CleanString.
func (as AddressList) String() string {
	return as.CleanString()
}

// Flatten returns the AddressList as a MailboxList. This returns a slice of
// Mailboxes. If the AddressList contains any groups, then the returned
// MailboxList will contain all the mailboxes within those groups.
func (as AddressList) Flatten() MailboxList {
	mbs := make(MailboxList, 0, len(as))
	for _, a := range as {
		switch v := a.(type) {
		case *Mailbox:
			mbs = append(mbs, v)
		case *AddrSpec:
			mb, _ := NewMailbox("", v, "")
			mbs = append(mbs, mb)
		case *Group:
			mbs = append(mbs, v.MailboxList()...)
		default:
			mb, _ := NewMailboxStr(
				a.DisplayName(),
				a.Address(),
				a.Comment(),
			)
			mbs = append(mbs, mb)
		}
	}

	return mbs
}

// ParseEmailAddress will parse any single email address. This could actually be
// multiple mailbox addresses if the email address is a group address. The
// object returned will either be a *Group or a *Mailbox.
//
// An error is returned if there's a problem during the parse. If the parse is
// partially successful, the address will be returned and a PartialParseError
// object will be returned.
func ParseEmailAddress(a string) (Address, error) {
	a = strings.TrimSpace(a)
	m, cs := rfc5322.MatchAddress([]byte(a))

	var address Address
	err := ApplyActions(m, &address)
	if err != nil {
		return nil, err
	}

	if len(cs) > 0 {
		return address, PartialParseError{string(cs)}
	}

	return address, nil
}

// ParseEmailAddressList will parse any list of addresses. Individual addresses
// may either be mailboxes or groups.
func ParseEmailAddressList(a string) (AddressList, error) {
	a = strings.TrimSpace(a)
	m, cs := rfc5322.MatchAddressList([]byte(a))

	var addresses AddressList
	err := ApplyActions(m, &addresses)
	if err != nil {
		return nil, err
	}

	if len(cs) > 0 {
		return addresses, PartialParseError{string(cs)}
	}

	return addresses, nil
}
