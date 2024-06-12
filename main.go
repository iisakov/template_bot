package main

import (
	"fmt"
	"log"
	"os"
	"party_bot/config"
	"party_bot/model"
	"party_bot/stl/tgstl"

	tg "github.com/iisakov/telegram-bot-api"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Не найден файл .env")
	}
	config.TOKEN = os.Getenv("BOT_TOKEN")
	config.PASS = os.Getenv("BOT_PASSWORD")

	config.BotStage, _ = config.BotStage.ReadBackup()

	for k, v := range map[string]model.Backupable{
		"Pair":      &config.PAIRS,
		"Moderator": &config.MODERATORS,
		"Customer":  &config.CUSTOMERS,
		"Admin":     &config.ADMINS,
		"Question":  &config.QUESTIONS} {

		model.ReadBackup(v, k)
	}

}

func main() {
	bot, err := tg.NewBotAPI(config.TOKEN)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				fmt.Println("Команда")
				tgstl.HandleMessageComand(*update.Message, bot, *config.BotStage)
			} else {
				fmt.Println("Плоский текст")
				tgstl.HandleMessagesText(*update.Message, bot, config.BotStage)
			}
		}

		if update.CallbackQuery != nil {
			fmt.Println("Нажатая инлайн кнопка")
			tgstl.HandleMessagesCallbackQuery(*update.CallbackQuery, bot, config.BotStage)
		}

		for k, v := range map[string]model.Backupable{
			"Pair":      &config.PAIRS,
			"Moderator": &config.MODERATORS,
			"Customer":  &config.CUSTOMERS,
			"Admin":     &config.ADMINS,
			"Question":  &config.QUESTIONS} {

			model.CreateBackup(&v, k)
		}
		config.BotStage.WriteBackup()
	}
}
