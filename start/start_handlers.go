package start

import (
	"log"

	"gdg-kld.ru/telecomma/bot"
	"gdg-kld.ru/telecomma/text"
	"gdg-kld.ru/telecomma/user"
	"gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

func RegisterHandlers(b *telebot.Bot, db *gorm.DB) {
	b.Handle("/start", func(m *telebot.Message) {
		_, err := b.Send(m.Sender, text.BotGreeting, &telebot.SendOptions{
			DisableWebPagePreview: true,
			ReplyMarkup: mainKeyboard(
				[]telebot.ReplyButton{
					{Text: text.RegisterMe},
				},
			),
		})
		if err != nil {
			log.Println(err)
		}
	})

	registration := register(b, db)
	b.Handle(text.RegisterMe, registration)
	b.Handle(text.SendBtn, registration)
	b.Handle(text.ReturnBtn, registration)

	b.Handle(text.DeleteBtn, func(m *telebot.Message) {
		err := db.Delete(user.User{ID: m.Sender.ID}).Error
		if err != nil {
			bot.Err(b, m.Sender, err)

			return
		}

		bot.Send(b, m.Sender, text.DeleteApproved, mainKeyboard(
			[]telebot.ReplyButton{
				{Text: text.ReturnBtn},
			},
		))
	})

	b.Handle(text.HowToUseBtn, textResponse(b, text.HowToUse))
	b.Handle(text.HowToStartBtn, textResponse(b, text.HowToStart))
}

func register(b *telebot.Bot, db *gorm.DB) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		sender := m.Sender

		newUser := user.NewUser(m.Sender)

		err := db.Save(&newUser).Error
		if err != nil {
			bot.Err(b, sender, err)

			return
		}

		bot.Send(b, sender, text.Registered(newUser.Name), mainKeyboard(
			[]telebot.ReplyButton{
				{Text: text.SendBtn},
				{Text: text.DeleteBtn},
			},
		))
	}
}

func mainKeyboard(secondRow []telebot.ReplyButton) *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		ResizeReplyKeyboard: true,
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				{Text: text.HowToStartBtn},
				{Text: text.HowToUseBtn},
			},
			secondRow,
		},
	}
}

func textResponse(b *telebot.Bot, text string) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		bot.Send(b, m.Sender, text)
	}
}
