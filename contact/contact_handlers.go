package contact

import (
	"fmt"

	"gdg-kld.ru/telecomma/log"
	"gdg-kld.ru/telecomma/user"
	"gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

func RegisterHandlers(b *telebot.Bot, db *gorm.DB) {
	b.Handle(telebot.OnText, func(m *telebot.Message) {
		// todo: assert registered
		// todo: sent start if not registered

		var users []user.User

		err := db.Find(&users, "name like '%"+m.Text+"%'").Error
		if err != nil {
			log.Err(b, m.Sender, err)
		}

		fmt.Printf("message: %s\n", m.Text)
		fmt.Printf("users: %v\n", users)

		message := "TODO: implement search"

		for _, u := range users {
			message += "\n" + u.Name
		}

		log.Send(b, m.Sender, message, &telebot.ReplyMarkup{})
	})
}
