package addr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObsAngleAddr(t *testing.T) {
	const str = "\"display name\" <@obs1.example.com,@obs2.example.com,@obs3.example.com:cur@example.com>"

	t.Parallel()

	email, err := ParseEmailAddress(str)
	assert.NoError(t, err)

	assert.Equal(t, "display name", email.DisplayName())
	assert.Equal(t, "cur@example.com", email.Address())
}

func TestObsMboxList(t *testing.T) {
	const str = "\"who\" <ok@example.com>, (obsolete comment with no address)"

	t.Parallel()

	ml, err := ParseEmailMailboxList(str)
	assert.NoError(t, err)

	assert.Equal(t, MailboxList{
		&Mailbox{
			displayName: "who",
			address:     &AddrSpec{"ok", "example.com", "ok@example.com"},
			comment:     "",
			original:    "\"who\" <ok@example.com>",
		},
	}, ml)
}

func TestObsAddrList(t *testing.T) {
	const str = "meh: \"who\" <ok@example.com>;, (obsolete comment with no address)"

	t.Parallel()

	al, err := ParseEmailAddressList(str)
	assert.NoError(t, err)

	assert.Equal(t, AddressList{
		&Group{
			displayName: "meh",
			mailboxList: MailboxList{
				&Mailbox{
					displayName: "who",
					address:     &AddrSpec{"ok", "example.com", "ok@example.com"},
					comment:     "",
					original:    "\"who\" <ok@example.com>",
				},
			},
			original: "meh: \"who\" <ok@example.com>;",
		},
	}, al)
}

func TestObsGroup(t *testing.T) {
	const str = "meh: (obsolete comments here), (obsolete comment there);"

	t.Parallel()

	g, err := ParseEmailGroup(str)
	assert.NoError(t, err)

	assert.Equal(t, &Group{
		displayName: "meh",
		mailboxList: MailboxList{},
		original:    str,
	}, g)
}
