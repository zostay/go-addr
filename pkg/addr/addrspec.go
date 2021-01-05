package addr

import (
	"fmt"

	"github.com/zostay/go-addr/pkg/format"
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
