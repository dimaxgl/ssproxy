package main

import (
	"flag"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

var (
	token = flag.String(`token`, ``, `telegram bot api token`)
)
// TODO implement administration from telegram bot
func main() {
	flag.Parse()

	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		log.Fatalln(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalln(err)
	}

	for u := range updates {
		log.Println(u)
	}
}
