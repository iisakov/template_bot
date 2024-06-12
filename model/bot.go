package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	tg "github.com/iisakov/telegram-bot-api"
)

type Backupable interface {
	CreateBackup()
	ReadBackup()
}

func CreateBackup(bv *Backupable, name string) {
	json, err := json.MarshalIndent(bv, "", "	")
	if err != nil {
		fmt.Println(err.Error())
	}
	os.WriteFile(fmt.Sprintf("backup/backup%s.json", name), json, 0666)
}

func ReadBackup(bv Backupable, name string) {
	if _, err := os.Stat(fmt.Sprintf("backup/backup%s.json", name)); errors.Is(err, os.ErrNotExist) {
		log.Println(err)
	} else {
		f, err := os.ReadFile(fmt.Sprintf("backup/backup%s.json", name))
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(f, bv)
		if err != nil {
			panic(err)
		}
	}
}

func RenderInlineMarkup(b *tg.BotAPI, q Question) (*tg.InlineKeyboardMarkup, bool) {
	if len(q.Answers) == 0 {
		return nil, false
	}
	var rows [][]tg.InlineKeyboardButton
	var buttons []tg.InlineKeyboardButton
	for i, a := range q.Answers {
		buttons = append(buttons, tg.NewInlineKeyboardButtonData(strconv.Itoa(i+1), a))
		if (i+1)%3 == 0 {
			rows = append(rows, tg.NewInlineKeyboardRow(buttons...))
			buttons = []tg.InlineKeyboardButton{}
		}
	}
	if len(buttons) != 0 {
		rows = append(rows, tg.NewInlineKeyboardRow(buttons...))
	}
	keyboardMarkup := tg.NewInlineKeyboardMarkup(rows...)

	return &keyboardMarkup, true
}
