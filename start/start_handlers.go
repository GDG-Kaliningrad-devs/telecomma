package start

import (
	"fmt"

	"gdg-kld.ru/telecomma/bot"
	"gdg-kld.ru/telecomma/text"
	"gdg-kld.ru/telecomma/user"
	"gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

func RegisterHandlers(b *telebot.Bot, db *gorm.DB) {
	b.Handle("/start", func(m *telebot.Message) {
		bot.Send(b, m.Sender, text.BotGreeting, &telebot.ReplyMarkup{
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
			bot.Err(b, m.Sender, err)

			return
		}

		bot.Send(b, m.Sender, text.DeleteApproved, &telebot.ReplyMarkup{
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

		newUser, err := user.NewUser(m.Sender)
		if err != nil {
			bot.Err(b, sender, err)

			return
		}

		result := db.Create(&newUser)
		if result.Error != nil {
			// todo: handle already registered 'UNIQUE constraint failed: users.id'
			bot.Err(b, sender, result.Error)

			return
		}

		fmt.Println("user created", newUser.Name)

		bot.Send(b, sender, text.Registered, &telebot.ReplyMarkup{
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
