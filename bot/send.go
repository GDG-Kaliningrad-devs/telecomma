package bot

import (
	"log"

	"gopkg.in/tucnak/telebot.v2"
)

func Respond(b *telebot.Bot, c *telebot.Callback, resp *telebot.CallbackResponse) {
	err := b.Respond(c, resp)
	if err != nil {
		log.Println(err)
	}
}

func Send(b *telebot.Bot, to telebot.Recipient, what interface{}, reply ...*telebot.ReplyMarkup) {
	var err error

	if len(reply) == 0 {
		_, err = b.Send(to, what)
	} else {
		_, err = b.Send(to, what, reply[0])
	}

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
