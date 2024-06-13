package tgstl

import (
	"fmt"
	"party_bot/config"
	"party_bot/model"
	"strings"

	tg "github.com/iisakov/telegram-bot-api"
)

func HandleMessagesText(um tg.Message, b *tg.BotAPI, s *tg.Stages) {
	var msg tg.MessageConfig
	var cId int64 = um.From.ID
	config.LMId = um.MessageID
	var r string
	var err error

	if u, ok := config.MODERATORS.FindUser(um.From.UserName); !ok { // Если сообщение пришло НЕ от модератора
		switch s.CurrentStageNum {
		case 0: // Настройки
			if config.PASS == um.Text { // Сообщение это пароль
				u := config.MODERATORS.AddUser(um)
				msg = tg.NewMessage(u.UserChat_id, "Привет. Теперь ты один из модераторов.")
			} else { // Сообщение не совпадает с паролем
				msg = tg.NewMessage(um.From.ID, "Привет. Нужен пароль. Его знает админ. Спроси у него.")
			}
			DeleteAllMessages(b, cId, int(config.LMId))
		case 1: // Регистрация
			if u, ok := config.CUSTOMERS.FindUser(um.From.UserName); !ok { // Ползователь ещё не зарегистрирован
				u := config.CUSTOMERS.AddUser(um)
				config.CUSTOMERS.UpdateAlias(u.Login, um.Text)
				msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Привет. Теперь твои сообщения будут подписаны %s. У тебя есть пара минтут на то чтобы изменить псевдоним.", u.Alias))
			} else {
				config.CUSTOMERS.UpdateAlias(u.Login, um.Text)
				msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Изменили твой псевдоним. Теперь твои сообщения будут подписаны %s. У тебя есть пара минтут на то чтобы изменить псевдоним.", u.Alias))
			}
			DeleteAllMessages(b, cId, int(config.LMId))
		case 2: // Ответы на вопросы
			fallthrough
		case 4: // Ответы на вопросы
			msg = tg.NewMessage(um.From.ID, "Вопросы со свободным ответом появятся позже, следите за обновлениями.")
			DeleteAllMessages(b, cId, int(config.LMId))
		case 3: // Общение в парах
			if up, ok := config.PAIRS.FindPartner(um.From.UserName); ok {
				u, _ = config.PAIRS.FindUser(um.From.UserName)
				msg = tg.NewMessage(up.UserChat_id, fmt.Sprintf("Сообщение от %s:\n%s", u.Alias, um.Text))
			}
		}
		fmt.Println(r, err)

	} else { // Если сообщение пришло от модератора
		// Проверяем всегда
		switch um.Text {
		case "Текущий этап":
			msg = tg.NewMessage(cId, fmt.Sprintf("Этап: %s", s.Value[s.CurrentStageNum].Name))

		case "Следующий этап":
			if s.Up() {
				msg = tg.NewMessage(cId, fmt.Sprintf("Перешёл на следующий этап: %s", s.Value[s.CurrentStageNum].Name))
			} else {
				msg = tg.NewMessage(cId, fmt.Sprintf("Этап: %s", s.Value[s.CurrentStageNum].Name))
			}
			msg.ReplyMarkup, _ = model.RenderKeyboardMarkup(b, config.COMANDS.GetComands(s.CurrentStageNum))

		case "Предыдущий этап":
			if s.Down() {
				msg = tg.NewMessage(cId, fmt.Sprintf("Перешёл на Предыдущий этап: %s", s.Value[s.CurrentStageNum].Name))
			} else {
				msg = tg.NewMessage(cId, fmt.Sprintf("Этап: %s", s.Value[s.CurrentStageNum].Name))
			}
			msg.ReplyMarkup, _ = model.RenderKeyboardMarkup(b, config.COMANDS.GetComands(s.CurrentStageNum))

		case "Показать пользователей":
			msg = tg.NewMessage(cId, config.CUSTOMERS.String())

		case "Показать пары":
			msg = tg.NewMessage(cId, config.PAIRS.String())
		case "Распределить пары":
			config.PAIRS = config.CUSTOMERS.DistributionPairs()
			msg = tg.NewMessage(cId, config.PAIRS.String())
		case "Номер пар":
			msg = tg.NewMessage(cId, config.PAIRS.NumPair())
		case "Отправить номер столика":
			recipients := ""
			for i, p := range config.PAIRS {
				msg = tg.NewMessage(p.Pair[0].UserChat_id, fmt.Sprintf("Пересядь, за столик - %d", i+1))
				b.Send(msg)
				msg = tg.NewMessage(p.Pair[1].UserChat_id, fmt.Sprintf("Пересядь, за столик - %d", i+1))
				b.Send(msg)
				recipients += fmt.Sprintf("%s (%s)\n", p.Pair[0].Login, p.Pair[0].Alias)
				recipients += fmt.Sprintf("%s (%s)\n", p.Pair[1].Login, p.Pair[1].Alias)
			}
			msg = tg.NewMessage(cId, fmt.Sprintf("Отправил столики.\n\nПо списку:\n%s", recipients))

		default:
			msg = tg.NewMessage(cId, "Что-то пошло не так")
		}

		// Проверяем в определённые этапы
		switch s.CurrentStageNum {
		case 0: // Настройки
			switch um.Text {
			case "Тест":
				msg = tg.NewMessage(cId, "Тестирую")
				msg.ReplyMarkup, _ = model.RenderKeyboardMarkup(b, config.COMANDS.GetComands(s.CurrentStageNum))
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
					msg = tg.NewMessage(u.UserChat_id, qText)
					msg.ReplyMarkup, _ = model.RenderInlineMarkup(b, config.QUESTIONS.GetQuestion())
					DeleteAllMessages(b, cId, int(config.LMId))
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
				break
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

			msg = tg.NewMessage(u.UserChat_id, fmt.Sprintf("Отправил задание:\n%s\n\nПо списку:\n%s", splitMessage[1], recipients))
		}

	}
	b.Send(msg)
}
