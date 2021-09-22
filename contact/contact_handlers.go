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
	b.Handle(&telebot.InlineButton{Unique: approve}, requestApprove(b, db))
	b.Handle(&telebot.InlineButton{Unique: decline}, requestDecline(b, db))
	b.Handle(telebot.OnCallback, requestContact(b, db))
}

func search(b *telebot.Bot, db *gorm.DB) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		// todo: assert 3 more symbols
		var users []user.User

		err := db.Find(&users, "name like '%"+m.Text+"%'").Error
		if err != nil {
			bot.Err(b, m.Sender, err)
		}

		if len(users) == 0 {
			bot.Send(b, m.Sender, text.SearchNotFound, &telebot.ReplyMarkup{})
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

		contactID, err := getContactID(b, c)
		if err != nil {
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
				// todo: implement "you ignored it, want to approve?"
				bot.Err(b, c.Sender, gorm.ErrNotImplemented)
			},
		)
		if exist {
			return
		}

		err = db.Create(user.NewContact(userID, contactID)).Error
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

func requestApprove(b *telebot.Bot, db *gorm.DB) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {
		// todo: assert registered
		// todo: approve request
		bot.Respond(b, c, &telebot.CallbackResponse{
			Text: text.ContactRequestSuccess,
		})
	}
}

func requestDecline(b *telebot.Bot, db *gorm.DB) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {
		// todo: assert registered
		// todo: decline request
		bot.Respond(b, c, &telebot.CallbackResponse{
			Text: text.ContactRequestDeclined,
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

func getContactID(b *telebot.Bot, c *telebot.Callback) (int, error) {
	data := strings.Split(c.Data, "|")

	if len(data) != 2 {
		fmt.Println(c.Data)

		bot.Respond(b, c, &telebot.CallbackResponse{
			Text: text.ContactRequestNotID,
		})

		return -1, errors.New("no id in data")
	}

	userID, err := strconv.Atoi(data[1])
	if err != nil {
		fmt.Println(c.Data)

		bot.Respond(b, c, &telebot.CallbackResponse{
			Text: text.ContactRequestNotID,
		})

		return 0, err
	}

	return userID, err
}
