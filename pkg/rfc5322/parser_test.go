package rfc5322

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zostay/go-addr/internal/rd"
)

func TestMatchAddressHappyMailbox(t *testing.T) {
	t.Parallel()

	mb := "\"ABC 123\" <abc213@example.com>"

	m, cs := MatchAddress([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TNameAddr, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchAddressHappyGroup(t *testing.T) {
	t.Parallel()

	mb := "Group: \"ABC 123\" <abc213@example.com>;"

	m, cs := MatchAddress([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TGroup, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchMailboxHappyNameAddr(t *testing.T) {
	t.Parallel()

	mb := "\"ABC 123\" <abc213@example.com>"

	m, cs := MatchMailbox([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TNameAddr, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchNameAddrHappy(t *testing.T) {
	t.Parallel()

	mb := "\"ABC 123\" <abc213@example.com>"

	m, cs := MatchNameAddr([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TNameAddr, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchAngleAddrHappyCur(t *testing.T) {
	t.Parallel()

	mb := "<abc213@example.com>"

	m, cs := MatchAngleAddr([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TAngleAddr, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchGroupHappy(t *testing.T) {
	t.Parallel()

	mb := "Group: \"ABC 123\" <abc213@example.com>;"

	m, cs := MatchAddress([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TGroup, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchDisplayNameHappy(t *testing.T) {
	t.Parallel()

	mb := "\"ABC 123\""

	m, cs := MatchDisplayName([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TDisplayName, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchMailboxListHappy(t *testing.T) {
	t.Parallel()

	mb := "\"ABC 123\" <abc213@example.com>, foo <bar@example.com>"

	m, cs := MatchMailboxList([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TMailboxList, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchAddressListHappy(t *testing.T) {
	t.Parallel()

	mb := "\"ABC 123\" <abc213@example.com>, Group: foo <bar@example.com>;"

	m, cs := MatchAddressList([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TAddressList, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchGroupListHappy(t *testing.T) {
	t.Parallel()

	mb := "\"ABC 123\" <abc213@example.com>, foo <bar@example.com>"

	m, cs := MatchGroupList([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TMailboxList, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchAddrSpecHappy(t *testing.T) {
	t.Parallel()

	mb := "abc213@example.com"

	m, cs := MatchAddrSpec([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TAddrSpec, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchLocalPartHappyDotAtom(t *testing.T) {
	t.Parallel()

	mb := "abc213"

	m, cs := MatchLocalPart([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TLiteral, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchDomainHappyDotAtom(t *testing.T) {
	t.Parallel()

	mb := "example.com"

	m, cs := MatchDomain([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TLiteral, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

// func TestMatchDomainLiteralHappy(t *testing.T) {
// 	t.Parallel()
//
// 	mb := "[127.0.0.1]"
//
// 	m, cs := MatchDomainLiteral([]byte(mb))
// 	assert.NotNil(t, m)
//
// 	assert.Empty(t, cs)
// 	assert.Equal(t, rd.TLiteral, m.Tag)
// 	assert.Equal(t, []byte(mb), m.Content)
// }

func TestMatchDTextHappy(t *testing.T) {
	t.Parallel()

	mb := "x"

	m, cs := MatchDText([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TNone, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchWordHappyAtom(t *testing.T) {
	t.Parallel()

	mb := "abc!#$"

	m, cs := MatchWord([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TAtom, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchPhraseHappyWords(t *testing.T) {
	t.Parallel()

	mb := "abc 123 !#$ \"foo bar baz\""

	m, cs := MatchPhrase([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TWords, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchATextHappy(t *testing.T) {
	t.Parallel()

	mb := "!"

	m, cs := MatchAText([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TLiteral, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchAtomHappy(t *testing.T) {
	t.Parallel()

	mb := "abc!#$"

	m, cs := MatchAtom([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TAtom, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchDotAtomTextHappy(t *testing.T) {
	t.Parallel()

	mb := "abc!#$.def&'*.123=?_"

	m, cs := MatchDotAtomText([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TLiteral, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchDotAtomHappy(t *testing.T) {
	t.Parallel()

	mb := "abc!#$.def&'*.123=?_"

	m, cs := MatchDotAtom([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TLiteral, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchFWSHappyCur(t *testing.T) {
	t.Parallel()

	mb := " \r\n\t"

	m, cs := MatchFWS([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TFWS, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchCTextHappyCur(t *testing.T) {
	t.Parallel()

	mb := "!"

	m, cs := MatchCText([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TCText, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchCContentHappyCText(t *testing.T) {
	t.Parallel()

	mb := "!"

	m, cs := MatchCContent([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TCText, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchCommentHappy(t *testing.T) {
	t.Parallel()

	mb := "(abc123)"

	m, cs := MatchComment([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TComment, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchCFWSHappyCFWSWithComment(t *testing.T) {
	t.Parallel()

	mb := "(abc) (123)"

	m, cs := MatchCFWS([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TCFWS, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchObsFWSHappy(t *testing.T) {
	t.Parallel()

	mb := " \r\n\t\r\n "

	m, cs := MatchObsFWS([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TObsFWS, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchQTextHappyPrintableASCII(t *testing.T) {
	t.Parallel()

	mb := "T"

	m, cs := MatchQText([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TLiteral, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchQContentHappyQText(t *testing.T) {
	t.Parallel()

	mb := "Z"

	m, cs := MatchQContent([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TLiteral, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

// func TestMatchQuotedStringHappy(t *testing.T) {
// 	t.Parallel()
//
// 	mb := "This ia a quoted string"
//
// 	m, cs := MatchQuotedString([]byte(mb))
// 	assert.NotNil(t, m)
//
// 	assert.Empty(t, cs)
// 	assert.Equal(t, rd.TLiteral, m.Tag)
// 	assert.Equal(t, []byte(mb), m.Content)
// }

func TestMatchObsNoWSCtlHappy(t *testing.T) {
	t.Parallel()

	mb := "\x0b"

	m, cs := MatchObsNoWSCtl([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TLiteral, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchObsCTextHappy(t *testing.T) {
	t.Parallel()

	mb := "\x0c"

	m, cs := MatchObsCText([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TLiteral, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchObsQTextHappy(t *testing.T) {
	t.Parallel()

	mb := "\x0e"

	m, cs := MatchObsDText([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TLiteral, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchObsQPHappy(t *testing.T) {
	t.Parallel()

	mb := "\\\n"

	m, cs := MatchObsQP([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TObsQP, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchObsPhraseHappy(t *testing.T) {
	t.Parallel()

	mb := "blah .blee. bloo"

	m, cs := MatchObsPhrase([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TObsPhrase, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchQuotedPairHappy(t *testing.T) {
	t.Parallel()

	mb := "\\n"

	m, cs := MatchQuotedPair([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TLiteral, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchObsAngleAddrHappy(t *testing.T) {
	t.Parallel()

	mb := "<@example.com,@example.com:abc123@example.com>"

	m, cs := MatchObsAngleAddr([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TObsAngleAddr, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchObsRouteHappy(t *testing.T) {
	t.Parallel()

	mb := "@example.com,@example.com:"

	m, cs := MatchObsRoute([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TObsRoute, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchObsDomainListHappy(t *testing.T) {
	t.Parallel()

	mb := "@example.com,@example.com"

	m, cs := MatchObsDomainList([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TObsDomainList, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchObsMboxListHappy(t *testing.T) {
	t.Parallel()

	mb := ",,,\"ABC 123\" <abc213@example.com>, foo <bar@example.com>"

	m, cs := MatchObsMboxList([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TObsMboxList, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchObsAddrListHappy(t *testing.T) {
	t.Parallel()

	mb := ",,, \"ABC 123\" <abc213@example.com>, Group: foo <bar@example.com>;"

	m, cs := MatchObsAddrList([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TObsAddrList, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchObsGroupListHappy(t *testing.T) {
	t.Parallel()

	mb := ",,, "

	m, cs := MatchObsGroupList([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TObsGroupList, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchobsLocalPartHappy(t *testing.T) {
	t.Parallel()

	mb := "abc213.456def"

	m, cs := MatchObsLocalPart([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, TObsLocalPart, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchObsDomainHappy(t *testing.T) {
	t.Parallel()

	mb := "example.com"

	m, cs := MatchDomain([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TLiteral, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}

func TestMatchObsDTextHappyObsNoWSCtl(t *testing.T) {
	t.Parallel()

	mb := "\x01"

	m, cs := MatchObsDText([]byte(mb))
	assert.NotNil(t, m)

	assert.Empty(t, cs)
	assert.Equal(t, rd.TLiteral, m.Tag)
	assert.Equal(t, []byte(mb), m.Content)
}
