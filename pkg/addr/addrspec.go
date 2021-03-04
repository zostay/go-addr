package addr

import (
	"fmt"
	"strings"

	"github.com/zostay/go-addr/pkg/format"
	"github.com/zostay/go-addr/pkg/rfc5322"
)

// AddrSpec is a concrete type for holding a single email address with no
// metadata. This is just the "local@domain" bit with no display name or comment
// information. It can also track the original string parsed to produce the
// object for the purpose of roundtripping.
type AddrSpec struct {
	localPart string
	domain    string
	original  string
}

// DisplayName always returns an empty string.
func (as *AddrSpec) DisplayName() string { return "" }

// Address is an alias for CleanString.
func (as *AddrSpec) Address() string { return as.String() }

// Comment always returns an empty string.
func (as *AddrSpec) Comment() string { return "" }

// LocalPart returns the part of the email address from before the at sign.
func (as *AddrSpec) LocalPart() string { return as.localPart }

// SetLocalPart sets the part of the email address before the at sign. This will
// also clear the original string if set.
func (as *AddrSpec) SetLocalPart(lp string) {
	as.localPart = lp
	as.original = ""
}

// Domain returns the part of the email address after the at sign.
func (as *AddrSpec) Domain() string { return as.domain }

// SetDomain sets the part of the email address after the at sign. This will
// also clear the original string if set.
func (as *AddrSpec) SetDomain(d string) {
	as.domain = d
	as.original = ""
}

// OriginalString returns the originally parsed string if that string is set.
// This is useful for roundtripping. However, if you are building a new email,
// it is best to use the CleanString. See
// https://tools.ietf.org/html/rfc5322#section-4
func (as *AddrSpec) OriginalString() string {
	return as.original
}

// CleanString will return a clean version of the email address suitable for use
// in new email messages.
func (as *AddrSpec) CleanString() string {
	return fmt.Sprintf("%s@%s",
		format.MaybeEscape(as.LocalPart(), false),
		as.Domain(),
	)
}

// String is an alias for CleanString.
func (as *AddrSpec) String() string { return as.OriginalString() }

// NewAddrSpec creates a new AddrSpec object from the given local part and
// domain.
func NewAddrSpec(localPart, domain string) *AddrSpec {
	return &AddrSpec{
		localPart: localPart,
		domain:    domain,
	}
}

// NewAddrSpecParsed creates a new AddrSpec object from the given local part and
// domain and stores an originally parsed string for roundtripping.
func NewAddrSpecParsed(lp, d, o string) *AddrSpec {
	return &AddrSpec{
		localPart: lp,
		domain:    d,
		original:  o,
	}
}

// ParseEmailAddrSpec parses the given email string for a bare email.
//
// On success it will return either the email address. If the email address
// successfully parses, but there is text remaining after the address, a
// PartialParseError will be returned as well.
//
// On failure, the object will be nil and an error will be returned.
func ParseEmailAddrSpec(a string) (*AddrSpec, error) {
	a = strings.TrimSpace(a)
	m, cs := rfc5322.MatchAddrSpec([]byte(a))

	var address *AddrSpec
	err := ApplyActions(m, &address)
	if err != nil {
		return nil, err
	}

	if len(cs) > 0 {
		return address, PartialParseError{string(cs)}
	}

	return address, nil
}
