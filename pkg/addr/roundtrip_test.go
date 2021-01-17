package addr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMailboxRoundtrip(t *testing.T) {
	mb, err := ParseEmailMailbox("\"Orson Scott Card\" <ender(weird comment placement)@example.com>")

	assert.NoError(t, err)

	assert.Equal(t, "\"Orson Scott Card\" <ender@example.com> (weird comment placement)", mb.String())
	assert.Equal(t, "\"Orson Scott Card\" <ender(weird comment placement)@example.com>", mb.OriginalString())
}
