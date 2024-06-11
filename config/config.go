package config

import tg "github.com/iisakov/telegram-bot-api"

var BotStage = tg.NewStages(
	0,
	func() (result []tg.Stage) {
		for i, s := range []string{"Настройка", "Регистрация", "Вопросы 1", "Общение в парах 1", "Настройка", "Вопросы 2", "Задания в парах"} {
			result = append(result, *tg.NewStage(s, uint16(i)))
		}
		return
	}()...)
