package addr_test

import (
	"errors"
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

// This example shows how you can recover from a partial parse, if you want.
func ExamplePartialParseError() {
	mb, err := addr.ParseEmailMailbox(
		"\"CS\" <charles.sheffield@example.com> and extra text",
	)

	var r string
	var ppe addr.PartialParseError
	if errors.As(err, &ppe) {
		r = ppe.Remainder
	} else if err != nil {
		panic(err)
	}

	fmt.Printf("Parsed: %s\n", mb)
	fmt.Printf("Remainder: %s\n", r)

	// Output:
	// Parsed: CS <charles.sheffield@example.com>
	// Remainder: and extra text
}
