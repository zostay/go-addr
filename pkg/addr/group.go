package addr

import (
	"strings"

	"github.com/zostay/go-addr/pkg/rfc5322"
)

type Group struct {
	displayName string
	mailboxList MailboxList
	original    string
}

func (g *Group) DisplayName() string      { return g.displayName }
func (g *Group) MailboxList() MailboxList { return g.mailboxList }
func (g *Group) Address() string          { return g.MailboxList().OriginalString() }
func (g *Group) Comment() string          { return "" }

func (g *Group) OriginalString() string {
	if g.original != "" {
		return g.original
	} else {
		return g.CleanString()
	}
}

func (g *Group) CleanString() string {
	var a strings.Builder
	a.WriteString(g.displayName)
	a.WriteString(": ")
	first := true
	for _, mb := range g.mailboxList {
		if !first {
			a.WriteString(", ")
		}
		a.WriteString(mb.CleanString())
		first = false
	}
	a.WriteString(";")
	return a.String()
}

func (g *Group) String() string { return g.OriginalString() }

type GroupList []*Group

func (gs GroupList) OriginalString() string {
	var a strings.Builder
	comma := false
	for _, g := range gs {
		if comma {
			a.WriteString(", ")
		}
		a.WriteString(g.OriginalString())
		comma = true
	}
	return a.String()
}

func (gs GroupList) CleanString() string {
	var a strings.Builder
	comma := false
	for _, g := range gs {
		if comma {
			a.WriteString(", ")
		}
		a.WriteString(g.CleanString())
		comma = true
	}
	return a.String()
}

func NewGroupParsed(dn string, l MailboxList, o string) *Group {
	return &Group{
		displayName: dn,
		mailboxList: l,
		original:    o,
	}
}

func ParseEmailGroup(a string) (*Group, error) {
	m, cs := rfc5322.MatchGroup([]byte(a))

	var group Group
	err := ApplyActions(m, &group)
	if err != nil {
		return nil, err
	}

	if len(cs) > 0 {
		return &group, ErrPartialParse
	}

	return &group, nil
}
