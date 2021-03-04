package addr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMailbox(t *testing.T) {
	mb, err := NewMailboxStr(
		"Peyton Randalf",
		"peyton.randalf@example.com",
		"Virginia House of Burgesses",
	)

	assert.NoError(t, err)

	assert.Equal(t, "Peyton Randalf", mb.DisplayName())
	assert.Equal(t, "peyton.randalf", mb.AddrSpec().LocalPart())
	assert.Equal(t, "example.com", mb.AddrSpec().Domain())
	assert.Equal(t, "Virginia House of Burgesses", mb.Comment())
	assert.Equal(t, "", mb.OriginalString())
	assert.Equal(t, "\"Peyton Randalf\" <peyton.randalf@example.com> (Virginia House of Burgesses)", mb.CleanString())
}

func TestUnicodeDecode(t *testing.T) {
	mb, err := ParseEmailMailbox("=?utf-8?q?St=C3=A9rl=C3=AF=C3=B1g?= <sterling@example.com> (=?utf-8?q?=C2=A1Hola,_se=C3=B1or!?=)")
	assert.NoError(t, err)

	assert.Equal(t, "Stérlïñg", mb.DisplayName())
	assert.Equal(t, "sterling@example.com", mb.Address())
	assert.Equal(t, "¡Hola, señor!", mb.Comment())
	assert.Equal(t, "=?utf-8?q?St=C3=A9rl=C3=AF=C3=B1g?= <sterling@example.com> (=?utf-8?q?=C2=A1Hola,_se=C3=B1or!?=)", mb.OriginalString())
}

func TestUnicodeEncode(t *testing.T) {
	mb, err := NewMailboxStr(
		"Stérlïñg",
		"sterling@example.com",
		"¡Hola, señor!",
	)

	assert.NoError(t, err)

	assert.Equal(t, "Stérlïñg", mb.DisplayName())
	assert.Equal(t, "sterling@example.com", mb.Address())
	assert.Equal(t, "¡Hola, señor!", mb.Comment())
	assert.Equal(t, "=?utf-8?q?St=C3=A9rl=C3=AF=C3=B1g?= <sterling@example.com> (=?utf-8?q?=C2=A1Hola,_se=C3=B1or!?=)", mb.CleanString())

	// Roundtrip in reverse?
	mb2, err := ParseEmailMailbox(mb.CleanString())
	assert.NoError(t, err)

	assert.Equal(t, "Stérlïñg", mb2.DisplayName())
	assert.Equal(t, "sterling@example.com", mb2.Address())
	assert.Equal(t, "¡Hola, señor!", mb2.Comment())
	assert.Equal(t, "=?utf-8?q?St=C3=A9rl=C3=AF=C3=B1g?= <sterling@example.com> (=?utf-8?q?=C2=A1Hola,_se=C3=B1or!?=)", mb2.OriginalString())
}

func TestOddity(t *testing.T) {
	mb, err := ParseEmailMailbox("\"Customer Support\"\n <support@example.com>")
	assert.NoError(t, err)

	assert.Equal(t, "Customer Support", mb.DisplayName())
	assert.Equal(t, "support@example.com", mb.Address())
}
