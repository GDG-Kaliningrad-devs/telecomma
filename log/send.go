package log

import (
	"log"

	"gopkg.in/tucnak/telebot.v2"
)

func Send(b *telebot.Bot, to telebot.Recipient, what interface{}, reply *telebot.ReplyMarkup) {
	_, err := b.Send(to, what, reply)
	if err != nil {
		log.Println(err)
	}
}

func Err(b *telebot.Bot, to telebot.Recipient, err error) {
	log.Println(err)

	_, err = b.Send(to, err.Error())
	if err != nil {
		log.Println(err)
	}
}
