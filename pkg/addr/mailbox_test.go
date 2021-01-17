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
