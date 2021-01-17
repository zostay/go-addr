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
	var testStrs = []string{
		"\"who\" <ok@example.com>, (obsolete comment with no address)",
		"\"who\" <ok@example.com>, (obsolete comment with no address), <another@example.com>",
	}

	var testResults = []MailboxList{
		{
			&Mailbox{
				displayName: "who",
				address:     &AddrSpec{"ok", "example.com", "ok@example.com"},
				comment:     "",
				original:    "\"who\" <ok@example.com>",
			},
		},
		{
			&Mailbox{
				displayName: "who",
				address:     &AddrSpec{"ok", "example.com", "ok@example.com"},
				comment:     "",
				original:    "\"who\" <ok@example.com>",
			},
			&Mailbox{
				displayName: "",
				address:     &AddrSpec{"another", "example.com", "another@example.com"},
				comment:     "",
				original:    "<another@example.com>",
			},
		},
	}

	t.Parallel()

	for i, str := range testStrs {
		expect := testResults[i]

		ml, err := ParseEmailMailboxList(str)
		assert.NoError(t, err)

		assert.Equal(t, expect, ml)
	}
}

func TestObsAddrList(t *testing.T) {
	var testStrs = []string{
		"meh: \"who\" <ok@example.com>;, (obsolete comment with no address)",
		"meh: \"who\" <ok@example.com>, (obsolete comment with no address);, <another@example.com>",
	}

	var testResults = []AddressList{
		{
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
		},
		{
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
				original: "meh: \"who\" <ok@example.com>, (obsolete comment with no address);",
			},
			&Mailbox{
				displayName: "",
				address:     &AddrSpec{"another", "example.com", "another@example.com"},
				comment:     "",
				original:    "<another@example.com>",
			},
		},
	}

	t.Parallel()

	for i, str := range testStrs {
		expect := testResults[i]

		al, err := ParseEmailAddressList(str)
		assert.NoError(t, err)

		assert.Equal(t, expect, al)
	}
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

func TestObsLocalPart(t *testing.T) {
	const str = "\"words\".in.\"email\".are.\"obsolete\"@example.com"

	t.Parallel()

	ml, err := ParseEmailAddrSpec(str)
	assert.NoError(t, err)

	assert.Equal(t, &AddrSpec{
		localPart: "words.in.email.are.obsolete",
		domain:    "example.com",
		original:  str,
	}, ml)
}

func TestObsDomain(t *testing.T) {
	const str = "okay@obs . example. (comments!?)com"

	t.Parallel()

	ml, err := ParseEmailAddrSpec(str)
	assert.NoError(t, err)

	assert.Equal(t, &AddrSpec{
		localPart: "okay",
		domain:    "obs.example.com",
		original:  str,
	}, ml)
}

func TestObsDText(t *testing.T) {
	const str = "okay@[\x01\\ \x02\\n\x03\x04]"

	t.Parallel()

	ml, err := ParseEmailAddrSpec(str)
	assert.NoError(t, err)

	assert.Equal(t, &AddrSpec{
		localPart: "okay",
		domain:    "[\x01\\ \x02\\n\x03\x04]",
		original:  str,
	}, ml)
}
