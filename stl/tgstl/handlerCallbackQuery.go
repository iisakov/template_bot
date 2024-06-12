package tgstl

import (
	"fmt"
	"party_bot/config"

	tg "github.com/iisakov/telegram-bot-api"
)

func HandleMessagesCallbackQuery(uc tg.CallbackQuery, b *tg.BotAPI, s *tg.Stages) {
	var msg tg.Chattable

	if u, ok := config.CUSTOMERS.FindUser(uc.From.UserName); ok {
		q := config.QUESTIONS.Questions[config.QUESTIONS.CurrentQuestionNum]
		switch s.CurrentStageNum {
		case 2:
			fallthrough
		case 4:
			switch {
			case q.FindOption("gender"):
				if uc.Data == "Юноша" {
					config.CUSTOMERS.SetGender(uc.From.UserName, 1)
				} else {
					config.CUSTOMERS.SetGender(uc.From.UserName, 0)
				}

				msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Отлично, %s. Мы записали твой ответ!", u.Alias))
			case q.FindOption("onlyOne"):
				config.CUSTOMERS.UpdateAnswer(u.Login, uc.Data, q.Answers)
				msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Отлично, %s. Мы записали твой ответ!\nТы можешь его изменить если хочешь.", u.Alias))
			default:
				config.CUSTOMERS.AddAnswer(u.Login, uc.Data)
				msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Отлично, %s. Мы записали твой ответ!\nТы можешь Выбрать несколько вариантов.", u.Alias))
			}
		}
		b.Send(msg)
	}
}
