package model

import (
	"fmt"
	"strings"

	tg "github.com/iisakov/telegram-bot-api"
)

type User struct {
	UserId        int64    `json:"userId"`
	UserChat_id   int64    `json:"userChatId"`
	Login         string   `json:"userLogin"`
	Alias         string   `json:"userAlias"`
	Answers       []string `json:"answers"`
	Gender        int      `json:"gender"`
	LastMessageId int      `json:"lastMessageId"`
}

func NewUser(m tg.Message) User {
	return User{UserId: m.From.ID,
		UserChat_id:   m.Chat.ID,
		Login:         m.From.UserName,
		Alias:         m.Text,
		Gender:        -1,
		LastMessageId: m.MessageID}
}

func (u User) String() string {
	return fmt.Sprintf(
		"id: %d\nchatId: %d\nlogin: %s\nalias: %s\nanswers:\n %s\ngender: %s\nlastMessageId: %d",
		u.UserId,
		u.UserChat_id,
		u.Login,
		u.Alias,
		strings.Join(u.Answers, "\n "),
		u.getGenderAsString(),
		u.LastMessageId,
	)
}

func (u User) getGenderAsString() string {
	switch u.Gender {
	case 1:
		return "male"
	case 0:
		return "female"
	default:
		return "unexpect"
	}
}

func (u *User) SetLastMessageId(messageid int) bool {
	u.LastMessageId = messageid
	return true
}
