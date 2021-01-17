package addr_test

import (
	"fmt"

	"github.com/zostay/go-addr/pkg/addr"
)

// Here is an example showing the difference between clean strings and original
// strings for full roundtripping on a Mailbox. ParseEmailAddress,
// ParseEmailAddrSpec, and ParseEmailGroup will work similarly.
func Example_mailboxRoundtripping() {
	mb, _ := addr.ParseEmailMailbox("\"Orson Scott Card\" <ender(weird comment placement)@example.com>")
	fmt.Println(mb)
	fmt.Println(mb.OriginalString())
	// Output:
	// "Orson Scott Card" <ender@example.com> (weird comment placement)
	// "Orson Scott Card" <ender(weird comment placement)@example.com>
}

// This example shows how the original is lost when parsing a mailbox list. The
// same will happen with an address list. Both the clean and original are the
// same in this case. Oddities within email addresses will be preserved, but
// other bits will not be.
func Example_mailboxListRoundtripping() {
	addresses := ", (weird stuff), \"J.R.R. Tolkein\" <j.r.r.tolkein@example.com>, \"C.S. Lewis\" <jack@example.com>, (wacky)"
	mbs, _ := addr.ParseEmailMailboxList(addresses)
	fmt.Println(mbs)
	fmt.Println(mbs.OriginalString())
	// Output:
	// "J.R.R. Tolkein" <j.r.r.tolkein@example.com>, "C.S. Lewis" <jack@example.com>
	// "J.R.R. Tolkein" <j.r.r.tolkein@example.com>, "C.S. Lewis" <jack@example.com>
}
