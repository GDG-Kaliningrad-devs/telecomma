package main

import (
	"os"
	"time"

	"gdg-kld.ru/telecomma/flag"
	"gopkg.in/tucnak/telebot.v2"
)

func settings() telebot.Settings {
	if flag.Debug() {
		return telebot.Settings{
			Token:  os.Getenv("BOT_TOKEN"),
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		}
	}

	panic("no prod")
}
