package contact

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gdg-kld.ru/telecomma/bot"
	"gdg-kld.ru/telecomma/flag"
	"gdg-kld.ru/telecomma/text"
	"gdg-kld.ru/telecomma/user"
	"gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

const (
	approve = "rapprove"
	decline = "rdecline"
)

func RegisterHandlers(b *telebot.Bot, db *gorm.DB) {
	b.Handle(telebot.OnText, search(b, db))
	b.Handle(&telebot.InlineButton{Unique: approve}, response(b, db, true))
	b.Handle(&telebot.InlineButton{Unique: decline}, response(b, db, false))
	b.Handle(telebot.OnCallback, requestContact(b, db))
}

func search(b *telebot.Bot, db *gorm.DB) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		if len(m.Text) < 4 {
			bot.Send(b, m.Sender, text.SearchErrMinSymbols, &telebot.ReplyMarkup{})
		}

		var users []user.User

		err := db.Find(&users, "name like '%"+m.Text+"%'").Error
		if err != nil {
			bot.Err(b, m.Sender, err)
		}

		if len(users) == 0 {
			bot.Send(b, m.Sender, text.SearchErrNoResults, &telebot.ReplyMarkup{})
		}

		if len(users) > 5 {
			bot.Send(b, m.Sender, text.SearchErrToManyResults, &telebot.ReplyMarkup{})
		}

		var keyboard [][]telebot.InlineButton

		for _, u := range users {
			id := strconv.Itoa(u.ID)
			btn := telebot.InlineButton{
				Unique: id,
				Text:   u.Name,
				Data:   id,
			}

			keyboard = append(keyboard, []telebot.InlineButton{btn})
		}

		bot.Send(b, m.Sender, text.SearchResult, &telebot.ReplyMarkup{
			InlineKeyboard: keyboard,
		})
	}
}

//nolint:funlen // refactor priority 3
func requestContact(b *telebot.Bot, db *gorm.DB) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {
		userID := c.Sender.ID

		contactID, has := getContactID(b, c)
		if !has {
			return
		}

		if flag.Debug() {
			contactID = userID
		}

		exist := assertNoContact(
			b, db, c.Sender, userID, contactID,
			func(contact user.Contact) {
				bot.Respond(b, c, &telebot.CallbackResponse{
					Text:      text.ContactRequestErrSent,
					ShowAlert: true,
				})
			},
		)
		if exist {
			return
		}

		exist = assertNoContact(
			b, db, c.Sender, contactID, userID,
			func(contact user.Contact) {
				bot.Respond(b, c, &telebot.CallbackResponse{
					Text:      text.ContactRequestErrIgnored,
					ShowAlert: true,
				})
			},
		)
		if exist {
			return
		}

		err := db.Create(user.NewContact(userID, contactID)).Error
		if err != nil {
			bot.Err(b, c.Sender, err)
		}

		data := strconv.Itoa(userID)

		bot.Send(b,
			&telebot.User{ID: contactID},
			text.ContactRequest(c.Sender.Username),
			&telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{{
					{
						Unique: approve,
						Text:   text.ContactRequestApproveBtn,
						Data:   data,
					},
					{
						Unique: decline,
						Text:   text.ContactRequestDeclineBtn,
						Data:   data,
					},
				}},
			},
		)

		bot.Respond(b, c, &telebot.CallbackResponse{
			Text: text.ContactRequestSend,
		})
	}
}

func response(b *telebot.Bot, db *gorm.DB, accepted bool) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {
		userID := c.Sender.ID

		contactID, has := getContactID(b, c)
		if !has {
			return
		}

		contact := user.Contact{SenderID: contactID, ContactID: userID}

		err := db.First(&contact).Error
		if err != nil {
			bot.Respond(b, c, &telebot.CallbackResponse{
				Text: err.Error(),
			})

			return
		}

		if contact.Response != user.None {
			bot.Respond(b, c, &telebot.CallbackResponse{
				Text: text.ContactResponseErrSent,
			})

			return
		}

		err = db.Save(contact.Respond(accepted)).Error
		if err != nil {
			bot.Respond(b, c, &telebot.CallbackResponse{
				Text: err.Error(),
			})

			return
		}

		receiver := &telebot.User{ID: contactID}

		if accepted {
			bot.Send(b, c.Sender, text.ContactRequestSuccess)
			bot.Send(b, receiver, text.ContactRequestSuccess)
		} else {
			bot.Send(b, receiver, text.ContactRequestDeclined)
		}

		bot.Respond(b, c, &telebot.CallbackResponse{
			Text: text.ActionDone,
		})
	}
}

func assertNoContact(
	b *telebot.Bot,
	db *gorm.DB,
	to telebot.Recipient,
	senderID, contactID int,
	onExist func(contact user.Contact),
) bool {
	//
	existedContact := user.Contact{SenderID: senderID, ContactID: contactID}

	err := db.First(&existedContact).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	if err != nil {
		bot.Err(b, to, err)

		return true
	}

	onExist(existedContact)

	return true
}

func getContactID(b *telebot.Bot, c *telebot.Callback) (int, bool) {
	data := strings.Split(c.Data, "|")

	userID, err := strconv.Atoi(data[len(data)-1])
	if err != nil {
		fmt.Println(c.Data)

		bot.Respond(b, c, &telebot.CallbackResponse{
			Text: text.ContactRequestNotID,
		})

		return 0, false
	}

	return userID, true
}
