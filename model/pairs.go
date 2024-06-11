package model

import "strings"

type Pairs []Pair

func (ps Pairs) unpack() (result []string) {
	for _, pair := range ps {
		result = append(result, pair.String())
	}
	return
}

func (ps Pairs) String() string {
	return strings.Join(ps.unpack(), "\n")
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

// pair is Backupable
func (ps Pairs) CreateBackup() {}
func (ps Pairs) ReadBackup()   {}
