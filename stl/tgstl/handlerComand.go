package tgstl

import (
	"party_bot/config"
	"party_bot/stl"

	tg "github.com/iisakov/telegram-bot-api"
)

func HandleMessageComand(um tg.Message, b *tg.BotAPI, s tg.Stages) {
	if u, ok := config.MODERATORS.FindUser(um.From.UserName); ok {
		switch s.CurrentStageNum {
		case 0: // Настройки
			u.SetLastMessageId(um.MessageID)
			msg := tg.NewMessage(u.UserChat_id, "Ты отправил команду.")
			DeleteMessegeByIds(b, u.UserChat_id, stl.MakeUIntSlice(u.LastMessageId-50, u.LastMessageId))
			b.Send(msg)
			switch um.Command() {
			default:
				msg = tg.NewMessage(u.UserChat_id, "Команда не найдена.")
				b.Send(msg)
			}
		}
	}
}
