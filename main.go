package main

import (
	"log"
	"os"
	"party_bot/config"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Не найден файл .env")
	}
	config.TOKEN = os.Getenv("BOT_TOKEN")

	//TODO model.ReadBackup()
}

func main() {

}
