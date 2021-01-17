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

// This example shows how the original
