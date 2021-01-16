# FFC 5322 Email Address Parser for Go

This aims at being a complete and correct RFC 5322 email address parser. It is
not aimed at being fast. I wish I could say the same for the built in `net/mail`
package, but as of this writing, the built-in library is incapable of parsing
obsolete forms (which make up much of the Internet). This is still unlikely to
be capable of parsing every address or list of addresses found in the wild. The
world of email is just a dangerous place.

This is written in pure Go using a hand-written recursive descent parser. This
is actually a port of a similar library I wrote for Raku.

## Synopsis

```go
package main

import (
    "fmt"

    "github.com/zostay/go-addr/pkg/addr"
)

func main() {
    addresses := "\"J.R.R. Tolkein\" <j.r.r.tolkein@example.com>, \"C.S. Lewis\" <jack@example.com>"
    as, err := addr.ParseEmailAddressList(addresses)
    if err != nil {
        panic(err)
    }

    for _, a := range as {
        fmt.Println("Name: " + a.DisplayName())
        fmt.Println("Addr: " + a.Address())
    }

    // Output: 
    // Name: J.R.R. Tolkein
    // Addr: j.r.r.tolkein@example.com
    // Name: C.S. Lewis
    // Addr: jack@example.com
}
```

## Description

This is intended to provide a more complete email address parsing apparatus that
I can trust and tweak to suit my needs.

## Email Address Formats Parsed

This library uses terminology taken from RFC 5322.

### Address

An email address is either a mailbox or a group. (See below.)

### Address List

An email address list is a list of of mailboxes or groups or mixture of both
separated by commas. 

For example:

```
"J.R.R. Tolkein" <j.r.r.tolkein@example.com>,
"C.S. Lewis" <jack@example.com>,
Distopian: "George Orwell" <bb.is.watching@example.com,
"Aldous Huxley" <soma@example.com>,
mockingjay@exampel.com;,
"Alexander Dumas" <ad@example.com>
```

This list contains three mailboxes and one group, which contains three more
mailboxes.

### Mailbox

An email mailbox is what most people think of as an email address. It's either
just a plain email address or it's a display name paired with an email address
in angle brackets or just an email address. Mailboxes may also contain comments,
which is just text inside of matched pairs of parentheses.

`"Douglas Adams" <da@example.com> (DON'T PANIC!)`

### Mailbox List

This is a list of mailboxes without groups separated by commas.

### Group

This is a single group mailbox, which is a name followed by a colon, followed by
zero or more mailboxes separated by commas and is ended by a semi-colon.

### AddrSpec

In RFC 5322, the email address with no decoration is called an addr-spec.
includes just the email address itself.

### Obsolete Forms

This parser also works for parsing obsolete forms that go all the way back to
RFC 822.

## Parsing Functions

These are all part of the `github.com/zostay/go-addr/pkg/addr` package. All of
these functions may return an error if there's a problem during the parse or if
parser can't find an address in the input.

If the input string starts with something that can be parsed as an email
address, but can't be completely parsed, then the part that can be parsed will
be returned and a `PartialParseError` will be returned as the error. This has a
`Remainder` field that can be used to retrieve the unparsed part:

```go
package main

import (
	"errors"
	"fmt"

	"github.com/zostay/go-addr/pkg/addr"
)

func main() {
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
    // Parsed: "CS" <charles.sheffield@example.com>
    // Remainder: and extra text
}
```

### addr.ParseEmailAddress

`func ParseEmailAddress(a string) (addr.Address, error)`

This parses a single email address. This can parse either a mailbox address or a
group list. The object returned will be either an `addr.Group` or an
`addr.Mailbox`.

### addr.ParseEmailAddressList

`func ParseEmailAddressList(a string) (addr.AddressList, error)`

An `addr.AddressList` is a slice of `addr.Address`. This will parse a list of
mailbox and/or group email addresses separated by commas.

### addr.ParseEmailMailbox

`func ParseEmailMailbox(a string) (*addr.Mailbox, error)`

This will parse a single mailbox address.

### addr.ParseMailMailboxList

`func ParseEmailMailboxList(a string) (addr.MailboxList, error)`

An `addr.MailboxList` is a slice of `*addr.Mailbox`. This will parse a list of
mailbox email addresses separated by commas.

### addr.ParseEmailGroup

`func ParseEmailGroup(a string) (*addr.Group, error)`

This will parse a single group address.

### addr.ParseEmailAddrSpec(a string) (*addr.AddrSpec, error)`

`func ParseEmailAddrSpec(a string) (*addr.AddrSpec, error)`

This will parse a single `addr-spec`.

## net/mail

If you want to convert mailbox email addresses from this library into those of
the `net/mail` package built-in to Go, you just need to do the following:

```
package main

import (
    "fmt"
    "net/mail"

    "github.com/zostay/go-addr/pkg/addr"
)

func main() {
    addrmb, _ := addr.ParseEmailMailbox("\"David Weber\" <honorh@example.com>")
    mailmb := mail.Address{addrmb.DisplayName(), addrmb.Address()}
    fmt.Println(mailmb)
    // Output: {David Weber honorh@example.com}
}
```

That's simple enough and there are enough other useful mail handling libraries
out there, I didn't bother providing a conversion helper. You can write your own
in a couple lines of Go.

## Copyright and License

Copyright 2020 Andrew Sterling Hanenkamp

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
