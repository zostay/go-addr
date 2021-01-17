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

func TestMailboxListRoundtripping(t *testing.T) {
	addresses := ", (weird stuff), \"J.R.R. Tolkein\" <j.r.r.tolkein@example.com>, \"C.S. Lewis\" <jack@example.com>, (wacky)"
	mbs, err := ParseEmailMailboxList(addresses)

	assert.NoError(t, err)

	assert.Equal(t,
		"\"J.R.R. Tolkein\" <j.r.r.tolkein@example.com>, \"C.S. Lewis\" <jack@example.com>",
		mbs.String(),
	)
	assert.Equal(t,
		"\"J.R.R. Tolkein\" <j.r.r.tolkein@example.com>, \"C.S. Lewis\" <jack@example.com>",
		mbs.OriginalString(),
	)
}
