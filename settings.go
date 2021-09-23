package main

import (
	"os"
	"time"

	"gopkg.in/tucnak/telebot.v2"
)

func settings() telebot.Settings {
	return telebot.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
}
