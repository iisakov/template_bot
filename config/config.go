package config

import (
	"party_bot/model"

	tg "github.com/iisakov/telegram-bot-api"
)

var TOKEN, PASS string

var QUESTIONS model.Questions

var ADMINS, MODERATORS, CUSTOMERS model.Users

var PAIRS model.Pairs

var BotStage = tg.NewStages(
	0,
	func() (result []tg.Stage) {
		for i, s := range []string{"Настройка", "Регистрация", "Вопросы 1", "Общение в парах 1", "Вопросы 2", "Задания в парах"} {
			result = append(result, *tg.NewStage(s, uint16(i)))
		}
		return
	}()...)
