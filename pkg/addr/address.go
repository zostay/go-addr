package addr

import (
	"strings"

	"github.com/zostay/go-addr/pkg/rfc5322"
)

type Address interface {
	DisplayName() string
	Address() string
	OriginalString() string
	CleanString() string
	Comment() string
}

type AddressList []Address

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
func ParseEmailAddress(a string) (Address, error) {
	m, cs := rfc5322.MatchAddress([]byte(a))

	var address Address
	err := ApplyActions(m, &address)
	if err != nil {
		return nil, err
	}

	if len(cs) > 0 {
		return address, ErrPartialParse
	}

	return address, nil
}

func ParseEmailAddressList(a string) (AddressList, error) {
	m, cs := rfc5322.MatchAddressList([]byte(a))

	var addresses AddressList
	err := ApplyActions(m, &addresses)
	if err != nil {
		return nil, err
	}

	if len(cs) > 0 {
		return addresses, ErrPartialParse
	}

	return addresses, nil
}
