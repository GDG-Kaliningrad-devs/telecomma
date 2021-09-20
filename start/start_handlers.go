package start

import (
	"fmt"
	"strings"

	"gdg-kld.ru/telecomma/log"
	"gdg-kld.ru/telecomma/text"
	"gdg-kld.ru/telecomma/user"
	"gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

func RegisterHandlers(b *telebot.Bot, db *gorm.DB) {
	b.Handle("/start", func(m *telebot.Message) {
		log.Send(b, m.Sender, text.BotGreeting, &telebot.ReplyMarkup{
			ResizeReplyKeyboard: true,
			ReplyKeyboard: [][]telebot.ReplyButton{
				{
					{
						Text: text.RegisterMe,
					},
				},
			},
		})
	})

	registration := register(b, db)
	b.Handle(text.RegisterMe, registration)
	b.Handle(text.ReturnBtn, registration)

	// todo: add delete dialog

	b.Handle(text.DeleteBtn, func(m *telebot.Message) {
		err := db.Delete(user.User{ID: m.Sender.ID}).Error
		if err != nil {
			log.Err(b, m.Sender, err)

			return
		}

		log.Send(b, m.Sender, text.DeleteApproved, &telebot.ReplyMarkup{
			ResizeReplyKeyboard: true,
			ReplyKeyboard: [][]telebot.ReplyButton{
				{
					{
						Text: text.ReturnBtn,
					},
				},
			},
		})
	})
}

func register(b *telebot.Bot, db *gorm.DB) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		sender := m.Sender
		// todo: validate name

		newUser := user.User{
			ID:   sender.ID,
			Name: strings.Join([]string{sender.FirstName, sender.LastName}, " "),
			// todo: get bio
		}

		result := db.Create(&newUser)
		if result.Error != nil {
			// todo: handle already registered 'UNIQUE constraint failed: users.id'
			log.Err(b, sender, result.Error)

			return
		}

		fmt.Println("user created", newUser.Name)

		log.Send(b, sender, text.Registered, &telebot.ReplyMarkup{
			ResizeReplyKeyboard: true,
			ReplyKeyboard: [][]telebot.ReplyButton{
				{
					{Text: text.HowToStartBtn},
				},
				{
					{Text: text.HowToUseBtn},
					{Text: text.DeleteBtn},
				},
			},
		})
	}
}
