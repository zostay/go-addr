package addr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const str = "Brotherhood: \"Winston Smith\" <winston.smith@recdep.minitrue> (Records Department), Julia <julia@ficdep.minitrue>;, user <user@oceania>"

func TestBasic(t *testing.T) {
	t.Parallel()

	email, err := ParseEmailAddressList(str)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(email))

	assert.Equal(t, "Brotherhood", email[0].DisplayName())
	assert.IsType(t, &Group{}, email[0])

	mbs := email[0].(*Group).MailboxList()
	assert.Equal(t, 2, len(mbs))

	assert.Equal(t, "Winston Smith", mbs[0].DisplayName())
	assert.Equal(t, "winston.smith", mbs[0].AddrSpec().LocalPart())
	assert.Equal(t, "recdep.minitrue", mbs[0].AddrSpec().Domain())
	assert.Equal(t, "Records Department", mbs[0].Comment())

	assert.Equal(t, "Julia", mbs[1].DisplayName())
	assert.Equal(t, "julia", mbs[1].AddrSpec().LocalPart())
	assert.Equal(t, "ficdep.minitrue", mbs[1].AddrSpec().Domain())
	assert.Empty(t, mbs[1].Comment())

	assert.IsType(t, &Mailbox{}, email[1])

	mb := email[1].(*Mailbox)
	assert.Equal(t, "user", mb.DisplayName())
	assert.Equal(t, "user", mb.AddrSpec().LocalPart())
	assert.Equal(t, "oceania", mb.AddrSpec().Domain())
	assert.Empty(t, mb.Comment())
}
