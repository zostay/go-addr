package addr

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zostay/go-addr/internal/rd"
	p "github.com/zostay/go-addr/pkg/rfc5322"
)

func TestApplyActionsTLiteralHappy(t *testing.T) {
	t.Parallel()

	mc := "testing123"
	m := &rd.Match{
		Tag:     rd.TLiteral,
		Content: []byte(mc),
	}

	var s string
	err := ApplyActions(m, &s)
	assert.NoError(t, err)

	assert.Equal(t, mc, m.Made)
	assert.Equal(t, mc, s)
}

func TestApplyActionsTNameAddrHappy(t *testing.T) {
	t.Parallel()

	email := "\"Zip\" <zip@example.com>"
	m, cs := p.MatchNameAddr([]byte(email))

	assert.NotNil(t, m)
	assert.Empty(t, cs)

	var mb *Mailbox
	err := ApplyActions(m, &mb)
	assert.NoError(t, err)

	assert.Equal(t, &Mailbox{
		displayName: "Zip",
		address:     &AddrSpec{"zip", "example.com", "zip@example.com"},
		comment:     "",
		original:    email,
	}, mb)
}

func TestApplyActionsTAngleAddrHappy(t *testing.T) {
	t.Parallel()

	aa := "<foo@example.com>"
	m, cs := p.MatchAngleAddr([]byte(aa))

	assert.NotNil(t, m)
	assert.Empty(t, cs)

	var as *AddrSpec
	err := ApplyActions(m, &as)
	assert.NoError(t, err)

	assert.Equal(t, &AddrSpec{
		localPart: "foo",
		domain:    "example.com",
		original:  "foo@example.com",
	}, as)
}

func TestApplyActionsTDisplayNameHappy(t *testing.T) {
	t.Parallel()

	dn := "\"Testing 123\""
	m, cs := p.MatchDisplayName([]byte(dn))

	assert.NotNil(t, m)
	assert.Empty(t, cs)

	var d string
	err := ApplyActions(m, &d)
	assert.NoError(t, err)

	assert.Equal(t, "Testing 123", d)
}

func TestApplyActionsTAddrSpecHappy(t *testing.T) {
	t.Parallel()

	as := "moomoo@example.com"
	m, cs := p.MatchAddrSpec([]byte(as))

	assert.NotNil(t, m)
	assert.Empty(t, cs)

	var a *AddrSpec
	err := ApplyActions(m, &a)
	assert.NoError(t, err)

	assert.Equal(t,
		&AddrSpec{"moomoo", "example.com", "moomoo@example.com"},
		a,
	)
}
