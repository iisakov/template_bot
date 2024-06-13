package model

import (
	"fmt"
	"strings"
)

type Pairs []Pair

func (ps Pairs) unpack() (result []string) {
	for _, pair := range ps {
		result = append(result, pair.String())
	}
	return
}

func (ps Pairs) String() string {
	return strings.Join(ps.unpack(), "\n\n+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+\n\n")
}

func (ps Pairs) FindUser(login string) (*User, bool) {
	for _, p := range ps {
		if u, ok := p.FindUser(login); ok {
			return u, true
		}
	}
	return nil, false
}

func (ps Pairs) FindPair(login string) (*Pair, bool) {
	for _, p := range ps {
		if _, ok := p.FindUser(login); ok {
			return &p, true
		}
	}
	return nil, false
}

func (ps Pairs) FindPartner(login string) (*User, bool) {
	if p, ok := ps.FindPair(login); ok {
		u := p.GetPartner(login)
		return &u, ok
	}
	return nil, false
}

func (ps Pairs) NumPair() (result string) {
	for i, p := range ps {
		result += fmt.Sprintf("num: %d\n pair: %s", i+1, p)
	}
	return
}

// pair is Backupable
func (ps Pairs) CreateBackup() {}
func (ps Pairs) ReadBackup()   {}
