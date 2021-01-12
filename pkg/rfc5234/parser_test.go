package rfc5234

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zostay/go-addr/internal/rd"
)

func TestMatchAlphaHappy(t *testing.T) {
	t.Parallel()

	m, cs := MatchAlpha([]byte("ABCxyz"))
	assert.NotNil(t, m)

	assert.Equal(t, []byte("BCxyz"), cs)
	assert.Equal(t, &rd.Match{
		Tag:     rd.TLiteral,
		Content: []byte{'A'},
	}, m)
}

func TestMatchAlphaSad(t *testing.T) {
	t.Parallel()

	m, cs := MatchAlpha([]byte("1BCxyz"))
	assert.Nil(t, m)
	assert.Nil(t, cs)
}

func TestMatchDigitHappy(t *testing.T) {
	t.Parallel()

	m, cs := MatchDigit([]byte("123456"))
	assert.NotNil(t, m)

	assert.Equal(t, []byte("23456"), cs)
	assert.Equal(t, &rd.Match{
		Tag:     rd.TLiteral,
		Content: []byte{'1'},
	}, m)
}

func TestMatchDigitSad(t *testing.T) {
	t.Parallel()

	m, cs := MatchDigit([]byte("ABCxyz"))
	assert.Nil(t, m)
	assert.Nil(t, cs)
}

func TestMatchCRHappy(t *testing.T) {
	t.Parallel()

	m, cs := MatchCR([]byte("\r"))
	assert.NotNil(t, m)

	assert.Equal(t, []byte{}, cs)
	assert.Equal(t, &rd.Match{
		Tag:     rd.TLiteral,
		Content: []byte{'\r'},
	}, m)
}

func TestMatchCRSad(t *testing.T) {
	t.Parallel()

	m, cs := MatchCR([]byte("\t"))
	assert.Nil(t, m)
	assert.Nil(t, cs)
}

func TestMatchLFHappy(t *testing.T) {
	t.Parallel()

	m, cs := MatchLF([]byte("\n"))
	assert.NotNil(t, m)

	assert.Equal(t, []byte{}, cs)
	assert.Equal(t, &rd.Match{
		Tag:     rd.TLiteral,
		Content: []byte{'\n'},
	}, m)
}

func TestMatchLFSad(t *testing.T) {
	t.Parallel()

	m, cs := MatchLF([]byte("\t"))
	assert.Nil(t, m)
	assert.Nil(t, cs)
}

func TestMatchCRLFHappy(t *testing.T) {
	t.Parallel()

	m, cs := MatchCRLF([]byte("\r\nfoo"))
	assert.NotNil(t, m)

	assert.Equal(t, []byte("foo"), cs)
	assert.Equal(t, &rd.Match{
		Tag:     rd.TLiteral,
		Content: []byte("\r\n"),
		Group:   map[string]*rd.Match{},
		Submatch: []*rd.Match{
			&rd.Match{
				Tag:     rd.TLiteral,
				Content: []byte{'\r'},
			},
			&rd.Match{
				Tag:     rd.TLiteral,
				Content: []byte{'\n'},
			},
		},
	}, m)
}

func TestMatchCRLFSad(t *testing.T) {
	t.Parallel()

	m, cs := MatchCRLF([]byte("\t\r\n"))
	assert.Nil(t, m)
	assert.Nil(t, cs)
}

func TestMatchDQuoteHappy(t *testing.T) {
	t.Parallel()

	m, cs := MatchDQuote([]byte("\""))
	assert.NotNil(t, m)

	assert.Equal(t, []byte{}, cs)
	assert.Equal(t, &rd.Match{
		Tag:     rd.TLiteral,
		Content: []byte{'"'},
	}, m)
}

func TestMatchDQuoteSad(t *testing.T) {
	t.Parallel()

	m, cs := MatchDQuote([]byte("'"))
	assert.Nil(t, m)
	assert.Nil(t, cs)
}

func TestMatchHTabHappy(t *testing.T) {
	t.Parallel()

	m, cs := MatchHTab([]byte("\t"))
	assert.NotNil(t, m)

	assert.Equal(t, []byte{}, cs)
	assert.Equal(t, &rd.Match{
		Tag:     rd.TLiteral,
		Content: []byte{'\t'},
	}, m)
}

func TestMatchHTabSad(t *testing.T) {
	t.Parallel()

	m, cs := MatchHTab([]byte("\n"))
	assert.Nil(t, m)
	assert.Nil(t, cs)
}

func TestMatchSPHappy(t *testing.T) {
	t.Parallel()

	m, cs := MatchSP([]byte(" "))
	assert.NotNil(t, m)

	assert.Equal(t, []byte{}, cs)
	assert.Equal(t, &rd.Match{
		Tag:     rd.TLiteral,
		Content: []byte{' '},
	}, m)
}

func TestMatchSPSad(t *testing.T) {
	t.Parallel()

	m, cs := MatchSP([]byte("\n"))
	assert.Nil(t, m)
	assert.Nil(t, cs)
}

func TestMatchWSPHappySP(t *testing.T) {
	t.Parallel()

	m, cs := MatchWSP([]byte(" "))
	assert.NotNil(t, m)

	assert.Equal(t, []byte{}, cs)
	assert.Equal(t, &rd.Match{
		Tag:     rd.TLiteral,
		Content: []byte{' '},
	}, m)
}

func TestMatchWSPHappyHTab(t *testing.T) {
	t.Parallel()

	m, cs := MatchWSP([]byte("\t"))
	assert.NotNil(t, m)

	assert.Equal(t, []byte{}, cs)
	assert.Equal(t, &rd.Match{
		Tag:     rd.TLiteral,
		Content: []byte{'\t'},
	}, m)
}

func TestMatchWSPSad(t *testing.T) {
	t.Parallel()

	m, cs := MatchWSP([]byte("\n"))
	assert.Nil(t, m)
	assert.Nil(t, cs)
}

func TestMatchVCharHappy(t *testing.T) {
	t.Parallel()

	m, cs := MatchVChar([]byte("ABCxyz"))
	assert.NotNil(t, m)

	assert.Equal(t, []byte("BCxyz"), cs)
	assert.Equal(t, &rd.Match{
		Tag:     rd.TLiteral,
		Content: []byte{'A'},
	}, m)
}

func TestMatchVCharSad(t *testing.T) {
	t.Parallel()

	m, cs := MatchVChar([]byte(" "))
	assert.Nil(t, m)
	assert.Nil(t, cs)
}
