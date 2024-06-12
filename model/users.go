package model

import (
	"party_bot/stl"
	"strings"

	tg "github.com/iisakov/telegram-bot-api"
)

type Users []User

func (us Users) unpack() (result []string) {
	for _, str := range us {
		result = append(result, str.String())
	}
	return
}

func (us Users) String() string {
	return strings.Join(us.unpack(), "\n\n")
}

func (us *Users) AddUser(m tg.Message) User {
	u := NewUser(m)
	*us = append(*us, u)
	return u
}

func (us Users) findUserIndex(login string) (int, bool) {
	for i, u := range us {
		if u.Login == login {
			return i, true
		}
	}
	return -1, false
}

func (us Users) FindUser(login string) (*User, bool) {
	for _, u := range us {
		if u.Login == login {
			return &u, true
		}
	}
	return nil, false
}

func (us Users) SetGender(login string, gender int) bool {
	i, ok := us.findUserIndex(login)
	if !ok {
		return false
	}
	us[i].Gender = gender
	return true
}

func (us Users) SetLastMessageId(login string, messageId int) bool {
	i, ok := us.findUserIndex(login)
	if !ok {
		return false
	}
	us[i].SetLastMessageId(messageId)
	return true
}

func (us Users) UpdateAlias(login string, newAlias string) bool {
	i, ok := us.findUserIndex(login)
	if !ok {
		return false
	}
	us[i].Alias = newAlias
	return true
}

func (us Users) AddAnswer(login string, newAnswer string) bool {
	i, ok := us.findUserIndex(login)
	if !ok {
		return false
	}
	us[i].Answers = append(us[i].Answers, newAnswer)
	return true
}

func (us Users) UpdateAnswer(login string, newAnswer string, answers []string) bool {
	i, ok := us.findUserIndex(login)
	if !ok {
		return false
	}
	for _, a := range answers {
		index, ok := stl.FindElementIndex(us[i].Answers, a)
		if !ok {
			continue
		}
		us[i].Answers = append(us[i].Answers[:index], us[i].Answers[index+1:]...)
	}
	us[i].Answers = append(us[i].Answers, newAnswer)
	return true
}

func (us Users) GetUsersByGender(gender int) (result Users) {
	for _, u := range us {
		if u.Gender == gender {
			result = append(result, u)
		}
	}

	return
}

func (us Users) splitByGender() (Users, Users) {
	var mUsers, fUsers Users
	for _, u := range us {
		if u.Gender == 0 {
			fUsers = append(fUsers, u)
		} else {
			mUsers = append(mUsers, u)
		}
	}
	return mUsers, fUsers
}

func (u User) countMatches(cU User) (result int) {
	var biggerS, smallerS []string
	if len(u.Answers) > len(cU.Answers) {
		biggerS, smallerS = u.Answers, cU.Answers
	} else {
		biggerS, smallerS = cU.Answers, u.Answers
	}

	for _, uS := range smallerS {
		for _, uB := range biggerS {
			if uS == uB {
				result += 1
			}
		}
	}
	return
}

func (us Users) DistributionPairs() Pairs {
	mUsers, fUsers := us.splitByGender()

	var biggerS, smallerS Users
	var maxMatches = 0
	var subResult = make(map[int][]Users)

	if len(mUsers) > len(fUsers) {
		biggerS, smallerS = mUsers, fUsers
	} else {
		biggerS, smallerS = fUsers, mUsers
	}

	for _, uS := range smallerS {
		for _, uB := range biggerS {
			numMatches := uS.countMatches(uB)
			if maxMatches < numMatches {
				maxMatches = numMatches
			}

			if _, ok := subResult[numMatches]; !ok {
				subResult[numMatches] = []Users{}
			}
			subResult[numMatches] = append(subResult[numMatches], Users{uS, uB})
		}
	}

	var result Pairs
	for i := maxMatches; i >= 0; i-- {
		for _, subPair := range subResult[i] {
			_, ok0 := result.FindUser(subPair[0].Login)
			_, ok1 := result.FindUser(subPair[1].Login)
			if ok0 || ok1 {
				continue
			}

			np := NewPair([2]User{subPair[0], subPair[1]}, i)
			result = append(result, np)
		}
	}
	return result
}

// pair is Backupable
func (us Users) CreateBackup() {}
func (us Users) ReadBackup()   {}
