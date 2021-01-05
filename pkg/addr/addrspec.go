package addr

import (
	"fmt"

	"github.com/zostay/go-addr/pkg/format"
	"github.com/zostay/go-addr/pkg/rfc5322"
)

type AddrSpec struct {
	localPart string
	domain    string
	original  string
}

func (as *AddrSpec) LocalPart() string { return as.localPart }
func (as *AddrSpec) Domain() string    { return as.domain }

func (as *AddrSpec) OriginalString() string {
	if as.original != "" {
		return as.original
	}
	return as.CleanString()
}

func (as *AddrSpec) CleanString() string {
	return fmt.Sprintf("%s@%s",
		format.MaybeEscape(as.LocalPart(), false),
		as.Domain(),
	)
}

func (as *AddrSpec) String() string { return as.OriginalString() }

func NewAddrSpecParsed(lp, d, o string) *AddrSpec {
	return &AddrSpec{
		localPart: lp,
		domain:    d,
		original:  o,
	}
}

func ParseEmailAddrSpec(a string) (*AddrSpec, error) {
	m, cs := rfc5322.MatchAddrSpec([]byte(a))

	var address AddrSpec
	err := ApplyActions(m, &address)
	if err != nil {
		return nil, err
	}

	if len(cs) > 0 {
		return &address, ErrPartialParse
	}

	return &address, nil
}
