package model

import (
	"fmt"
)

type Pair struct {
	Pair       [2]User `json:"pair"`
	NumMatches int     `json:"numMatches"`
}

func NewPair(us [2]User, nm int) Pair {
	return Pair{Pair: us, NumMatches: nm}
}

func (p Pair) String() string {
	return fmt.Sprintf(
		"pair:\n %s\nnumMatches: %d",
		fmt.Sprintf("user1: %s:\n user2: %s",
			p.Pair[0],
			p.Pair[1],
		),
		p.NumMatches,
	)
}

func (p Pair) GetUsers() [2]string {
	return [2]string{p.Pair[0].Login, p.Pair[1].Login}
}

func (p Pair) FindUser(login string) (*User, bool) {
	var users Users = p.Pair[:]
	return users.FindUser(login)
}

func (p Pair) GetPartner(login string) User {
	var users Users = p.Pair[:]
	i, _ := users.findUserIndex(login)
	return p.Pair[1-i]
}
