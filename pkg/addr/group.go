package addr

import (
	"strings"
)

type Group struct {
	displayName string
	mailboxList MailboxList
	original    string
}

func (g *Group) DisplayName() string      { return g.displayName }
func (g *Group) MailboxList() MailboxList { return g.mailboxList }

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
	}
	a.WriteString(";")
	return a.String()
}

func (g *Group) String() string { return g.OriginalString() }

type GroupList []*Group

func NewGroupParsed(dn string, l MailboxList, o string) *Group {
	return &Group{
		displayName: dn,
		mailboxList: l,
		original:    o,
	}
}
