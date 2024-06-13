package tgstl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"party_bot/stl"

	tg "github.com/iisakov/telegram-bot-api"
)

var tgHostUrl = "https://api.telegram.org/bot"

func DeleteMessegeById(b *tg.BotAPI, chat_id int64, message_id int) (resp *http.Response, err error) {
	type DM struct {
		Chat_id    int64 `json:"chat_id"`
		Message_id int   `json:"message_id"`
	}

	data, err := json.Marshal(DM{Chat_id: chat_id, Message_id: message_id})
	if err != nil {
		return
	}

	r := bytes.NewReader(data)
	resp, err = http.Post(b.Token+"/deleteMessage", "application/json", r)
	if err != nil {
		return
	}
	return
}

func DeleteMessegeByIds(b *tg.BotAPI, chat_id int64, message_ids []int) (result string, err error) {
	type DM struct {
		Chat_id     int64 `json:"chat_id"`
		Message_ids []int `json:"message_ids"`
	}
	data, err := json.Marshal(DM{Chat_id: chat_id, Message_ids: message_ids})
	if err != nil {
		return
	}

	r := bytes.NewReader(data)
	resp, err := http.Post(tgHostUrl+b.Token+"/deleteMessages", "application/json", r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	result = string(bodyBytes)

	return
}

func DeleteAllMessages(b *tg.BotAPI, chat_id int64, message_id int) {
	result, err := DeleteMessegeByIds(b, chat_id, stl.MakeUIntSlice(message_id-message_id%100, message_id))
	fmt.Println(result, err, message_id-message_id%100, message_id)
	for i := int(message_id / 100); i > int(message_id/100)-3; i-- {
		result, err = DeleteMessegeByIds(b, chat_id, stl.MakeUIntSlice(i*100-99, i*100))
		fmt.Println(result, err, i*100-99, i*100)
	}
}
