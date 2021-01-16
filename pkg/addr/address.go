// Package addr is the primary entry point for working with email addresses.
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

// ParseEmailAddress will parse any single email address. This could actually be
// multiple mailbox addresses if the email address is a group address. The
// object returned will either be a *Group or a *Mailbox.
//
// An error is returned if there's a problem during the parse. If the parse is
// partially successful, the address will be returned and a PartialParseError
// object will be returned.
func ParseEmailAddress(a string) (Address, error) {
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
