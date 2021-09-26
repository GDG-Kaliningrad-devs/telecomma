package main

import (
	"log"

	"gdg-kld.ru/telecomma/bot"
	"gdg-kld.ru/telecomma/contact"
	"gdg-kld.ru/telecomma/start"
	"gdg-kld.ru/telecomma/user"
	"github.com/joho/godotenv"
	"gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not loaded")
	}

	db, err := gorm.Open(sqlite.Open("telecomma.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)

		return
	}

	err = db.AutoMigrate(&user.User{}, &user.Contact{})
	if err != nil {
		log.Fatal(err)

		return
	}

	b, err := telebot.NewBot(settings())
	if err != nil {
		log.Fatal(err)

		return
	}

	start.RegisterHandlers(b, db)
	contact.RegisterHandlers(b, db)
	b.Handle("/version", func(m *telebot.Message) {
		bot.Send(b, m.Sender, "1.2.1")
	})

	b.Start()
}
