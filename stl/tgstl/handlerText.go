package tgstl

import (
	"fmt"
	"party_bot/config"
	"party_bot/model"
	"party_bot/stl"
	"strings"

	tg "github.com/iisakov/telegram-bot-api"
)

func HandleMessagesText(um tg.Message, b *tg.BotAPI, s *tg.Stages) {
	var msg tg.Chattable
	var cId int64 = um.From.ID
	var mId int = um.MessageID

	if u, ok := config.MODERATORS.FindUser(um.From.UserName); !ok { // Если сообщение пришло НЕ от модератора
		switch s.CurrentStageNum {
		case 0: // Настройки
			if config.PASS == um.Text { // Сообщение это пароль
				u := config.MODERATORS.AddUser(um)
				msg = tg.NewMessage(u.UserChat_id, "Привет. Теперь ты один из модераторов.")
			} else { // Сообщение не совпадает с паролем
				msg = tg.NewMessage(um.From.ID, "Привет. Нужен пароль. Его знает админ. Спроси у него.")
			}
			DeleteMessegeByIds(b, cId, stl.MakeUIntSlice(mId-50, mId))
		case 1:
			if u, ok := config.CUSTOMERS.FindUser(um.From.UserName); !ok {
				u := config.CUSTOMERS.AddUser(um)
				config.CUSTOMERS.UpdateAlias(u.Login, um.Text)
				msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Привет. Теперь твои сообщения будут подписаны %s. У тебя есть пара минтут на то чтобы изменить псевдоним.", u.Alias))
			} else {
				config.CUSTOMERS.UpdateAlias(u.Login, um.Text)
				msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Изменили твой псевдоним. Теперь твои сообщения будут подписаны %s. У тебя есть пара минтут на то чтобы изменить псевдоним.", u.Alias))
			}
			DeleteMessegeByIds(b, cId, stl.MakeUIntSlice(mId-50, mId))
		case 2:
			fallthrough
		case 4:
			msg = tg.NewMessage(um.From.ID, "Вопросы со свободным ответом появятся позже, следите за обновлениями.")
			DeleteMessegeByIds(b, cId, stl.MakeUIntSlice(mId-50, mId))
		case 3:
			if u, ok := config.PAIRS.FindPartner(um.From.UserName); ok {
				msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Сообщение от %s:\n%s", u.Alias, um.Text))
			}

		}

	} else { // Если сообщение пришло от модератора
		// Проверяем всегда
		switch um.Text {
		case "Текущий этап":
			msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Этап: %s", s.Value[s.CurrentStageNum].Name))

		case "Следующий этап":
			if s.Up() {
				msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Перешёл на следующий этап: %s", s.Value[s.CurrentStageNum].Name))
			} else {
				msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Этап: %s", s.Value[s.CurrentStageNum].Name))
			}

		case "Предыдущий этап":
			if s.Down() {
				msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Перешёл на Предыдущий этап: %s", s.Value[s.CurrentStageNum].Name))
			} else {
				msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Этап: %s", s.Value[s.CurrentStageNum].Name))
			}

		case "Показать пользователей":
			msg = tg.NewMessage(u.UserChat_id, config.CUSTOMERS.String())

		case "Показать пары":
			msg = tg.NewMessage(u.UserChat_id, config.PAIRS.String())
		case "Распредилить пары":
			config.PAIRS = config.CUSTOMERS.DistributionPairs()
			msg = tg.NewMessage(u.UserChat_id, config.PAIRS.String())

		default:
			msg = tg.NewMessage(um.From.ID, "Что-то пошло не так")
		}

		// Проверяем в определённые этапы
		switch s.CurrentStageNum {
		case 0: // Настройки
			switch um.Text {
			case "Тест":
				msg = tg.NewMessage(u.UserChat_id, "Тестирую")
			}
		case 2:
			fallthrough
		case 4:
			switch um.Text {
			case "Отправить вопрос":
				recipients := ""
				for _, u := range config.CUSTOMERS {
					qText := fmt.Sprintf(
						"Вопрос:\n%s\nВарианты ответа:\n%s",
						config.QUESTIONS.GetQuestion().Text,
						config.QUESTIONS.GetQuestion().GetNumberedAnswers(),
					)
					var msg tg.MessageConfig = tg.NewMessage(u.UserChat_id, qText)
					msg.ReplyMarkup, _ = model.RenderInlineMarkup(b, config.QUESTIONS.GetQuestion())
					b.Send(msg)
					recipients += fmt.Sprintf("%s (%s)\n", u.Login, u.Alias)
				}
				msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Отправил вопрос:\n%s\n\nПо списку:\n%s", config.QUESTIONS.GetQuestion().Text, recipients))
			case "Текущий вопрос":
				q := config.QUESTIONS.Questions[config.QUESTIONS.CurrentQuestionNum]
				msg = tg.NewMessage(u.UserChat_id, q.String())
			case "Следующий вопрос":
				q, _ := config.QUESTIONS.Next()
				msg = tg.NewMessage(u.UserChat_id, q.String())
			case "Предыдущий вопрос":
				q, _ := config.QUESTIONS.Back()
				msg = tg.NewMessage(u.UserChat_id, q.String())
			}
		case 5:
			recipients := ""
			splitMessage := strings.Split(um.Text, " : ")
			if len(splitMessage) != 3 {
				var msg tg.MessageConfig = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Сообщение %s должно состоять из трёх частей, разделённых ' : '.\nНапример: userName : Задача : pair", um.Text))
				b.Send(msg)
			}
			if splitMessage[2] == "pair" {
				if p, ok := config.PAIRS.FindPair(splitMessage[0]); ok {
					for _, u := range p.Pair {
						var msg tg.MessageConfig = tg.NewMessage(u.UserChat_id, splitMessage[1])
						b.Send(msg)
						recipients += fmt.Sprintf("%s (%s)\n", u.Login, u.Alias)
					}
				}
			}
			if splitMessage[2] == "user" {
				if u, ok := config.CUSTOMERS.FindUser(splitMessage[0]); ok {
					var msg tg.MessageConfig = tg.NewMessage(u.UserChat_id, splitMessage[1])
					b.Send(msg)
					recipients += fmt.Sprintf("%s (%s)\n", u.Login, u.Alias)
				}
			}

			msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Отправил задание:\n%s\n\nПо списку:\n%s", config.QUESTIONS.GetQuestion().Text, recipients))
		}
		u.SetLastMessageId(mId)
		DeleteMessegeByIds(b, cId, stl.MakeUIntSlice(mId-50, mId))
	}
	b.Send(msg)
}