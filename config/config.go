package config

import (
	"party_bot/model"

	tg "github.com/iisakov/telegram-bot-api"
)

var TOKEN, PASS string
var LMId int

var QUESTIONS model.Questions

var ADMINS, MODERATORS, CUSTOMERS model.Users

var PAIRS model.Pairs

var COMANDS model.Comands = model.Comands{
	0: [][]model.Comand{
		{model.Comand{Text: "1"}, model.Comand{Text: "2"}, model.Comand{Text: "3"}},
		{model.Comand{Text: "4"}, model.Comand{Text: "5"}, model.Comand{Text: "6"}},
		{model.Comand{Text: "7"}, model.Comand{Text: "8"}, model.Comand{Text: "9"}}}}

var BotStage = tg.NewStages(
	0,
	func() (result []tg.Stage) {
		for i, s := range []string{"Настройка", "Регистрация", "Вопросы 1", "Общение в парах", "Вопросы 2", "Задания в парах"} {
			result = append(result, *tg.NewStage(s, uint16(i)))
		}
		return
	}()...)
