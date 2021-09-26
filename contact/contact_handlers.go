package contact

import (
	"errors"
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
	accept             = "rapprove"
	decline            = "rdecline"
	fakeAccept         = "rfake_accept"
	senderInterested   = "rset_sender_interested"
	receiverInterested = "rset_receiver_interested"
)

func RegisterHandlers(b *telebot.Bot, db *gorm.DB) {
	b.Handle(telebot.OnText, search(b, db))
	b.Handle(&telebot.InlineButton{Unique: accept}, response(b, db, user.Accepted))
	b.Handle(&telebot.InlineButton{Unique: decline}, response(b, db, user.Declined))
	b.Handle(&telebot.InlineButton{Unique: fakeAccept}, response(b, db, user.FakeAccepted))
	b.Handle(&telebot.InlineButton{Unique: senderInterested}, setInterested(b, db, true))
	b.Handle(&telebot.InlineButton{Unique: receiverInterested}, setInterested(b, db, false))
	b.Handle(telebot.OnCallback, requestContact(b, db))
	b.Handle("/top", top(b, db))
	b.Handle("/admin_notify_top_"+flag.AdminPass(), notifyAboutTop(b, db))
}

func search(b *telebot.Bot, db *gorm.DB) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		if len(m.Text) < 4 {
			bot.Send(b, m.Sender, text.SearchErrMinSymbols, &telebot.ReplyMarkup{})

			return
		}

		var users []user.User

		err := db.Find(
			&users,
			"id <> ? AND name like '%"+m.Text+"%'",
			m.Sender.ID,
		).Error
		if err != nil {
			bot.Err(b, m.Sender, err)
		}

		if len(users) == 0 {
			bot.Send(b, m.Sender, text.SearchErrNoResults, &telebot.ReplyMarkup{})

			return
		}

		if len(users) > 5 {
			bot.Send(b, m.Sender, text.SearchErrToManyResults, &telebot.ReplyMarkup{})

			return
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
			text.ContactRequest(user.Name(c.Sender), c.Sender.Username),
			&telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						{
							Unique: accept,
							Text:   text.ContactRequestApproveBtn,
							Data:   data,
						},
						{
							Unique: decline,
							Text:   text.ContactRequestDeclineBtn,
							Data:   data,
						},
					},
					{
						{
							Unique: fakeAccept,
							Text:   text.ContactRequestFakeAcceptBtn,
							Data:   data,
						},
					},
				},
			},
		)

		bot.Respond(b, c, &telebot.CallbackResponse{
			Text: text.ContactRequestSend,
		})
	}
}

// nolint:funlen // refactor priority 3
func response(b *telebot.Bot, db *gorm.DB, status user.Response) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {
		currentUser := user.NewUser(c.Sender)
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

		err = db.Save(contact.Respond(status)).Error
		if err != nil {
			bot.Respond(b, c, &telebot.CallbackResponse{
				Text: err.Error(),
			})

			return
		}

		receiver := &telebot.User{ID: contactID}

		switch status {
		case user.Accepted, user.FakeAccepted:
			contactUser := getUserFromRequest(contactID, c.Message.Text)

			success := text.ContactRequestSuccess(
				currentUser.Name,
				currentUser.UserName,
				contactUser.Name,
				contactUser.UserName,
			)

			bot.Send(b, receiver, success, interestedBtn(userID, true, false))

			if status == user.Accepted {
				bot.Send(b, c.Sender, success, interestedBtn(contactID, false, false))
			}

		case user.Declined:
			bot.Send(b, receiver, text.ContactRequestDeclined(currentUser.Name))

		case user.None:
		}

		bot.Respond(b, c, &telebot.CallbackResponse{
			Text: text.ActionDone,
		})
	}
}

func top(b *telebot.Bot, db *gorm.DB) func(m *telebot.Message) {
	calculator := statCalculator{}

	return func(m *telebot.Message) {
		stats, t, err := calculator.top(db, m.ID)
		if err != nil {
			bot.Err(b, m.Sender, err)

			return
		}

		bot.Send(b, m.Sender, text.Top(stats, t))
	}
}

func notifyAboutTop(b *telebot.Bot, db *gorm.DB) func(m *telebot.Message) {
	calculator := statCalculator{}

	return func(m *telebot.Message) {
		ids, err := calculator.ids(db)
		if err != nil {
			bot.Err(b, m.Sender, err)

			return
		}

		for _, id := range ids {
			bot.Send(b, &telebot.User{ID: id}, text.LookAtTop)
		}
	}
}

func interestedBtn(dataID int, sender, set bool) *telebot.ReplyMarkup {
	unique := senderInterested
	if !sender {
		unique = receiverInterested
	}

	btnText := text.ContactSetImportantBtn
	if set {
		btnText = "🔥" + btnText
	}

	return &telebot.ReplyMarkup{InlineKeyboard: [][]telebot.InlineButton{{
		{
			Unique: unique,
			Text:   btnText,
			Data:   strconv.Itoa(dataID),
		},
	}}}
}

func setInterested(b *telebot.Bot, db *gorm.DB, sender bool) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {
		userID := c.Sender.ID

		contactID, has := getContactID(b, c)
		if !has {
			return
		}

		contact := user.Contact{SenderID: userID, ContactID: contactID}
		if !sender {
			contact = user.Contact{SenderID: contactID, ContactID: userID}
		}

		err := db.First(&contact).Error
		if err != nil {
			bot.Respond(b, c, &telebot.CallbackResponse{
				Text: err.Error(),
			})

			return
		}

		contact = contact.ToggleInterested(sender)

		err = db.Save(contact).Error
		if err != nil {
			bot.Respond(b, c, &telebot.CallbackResponse{
				Text: err.Error(),
			})

			return
		}

		_, err = b.EditReplyMarkup(c.Message, interestedBtn(contactID, sender, contact.Interested(sender)))
		if err != nil {
			bot.Err(b, c.Sender, err)
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
		bot.Respond(b, c, &telebot.CallbackResponse{
			Text: text.ContactRequestWrongData,
		})

		return 0, false
	}

	return userID, true
}

func getUserFromRequest(contactID int, s string) *user.User {
	split := strings.Split(s, ",")

	name := split[0]
	username := strings.Split(strings.TrimSpace(split[1]), " ")[0]

	return &user.User{
		ID:       contactID,
		Name:     name,
		UserName: strings.TrimPrefix(username, "@"),
	}
}
