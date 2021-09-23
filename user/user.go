package user

import (
	"time"

	"gopkg.in/tucnak/telebot.v2"
)

type User struct {
	ID       int // telegram user id
	Name     string
	UserName string
}

func NewUser(user *telebot.User) User {
	return User{
		ID:       user.ID,
		Name:     Name(user),
		UserName: user.Username,
	}
}

type Contact struct { //nolint:maligned // dummy optimization
	SenderID          int `gorm:"primaryKey"`
	SenderInterested  bool
	ContactID         int `gorm:"primaryKey"`
	ContactInterested bool
	Response          Response
	RequestTime       time.Time
	ResponseTime      *time.Time
}

func NewContact(firstID int, secondID int) Contact {
	return Contact{
		SenderID:    firstID,
		ContactID:   secondID,
		Response:    None,
		RequestTime: time.Now(),
	}
}

func (c Contact) Respond(status Response) Contact {
	c.Response = status

	now := time.Now()

	c.ResponseTime = &now

	return c
}

func (c Contact) ToggleInterested(sender bool) Contact {
	if sender {
		c.SenderInterested = !c.SenderInterested
	} else {
		c.ContactInterested = !c.ContactInterested
	}

	return c
}

func (c Contact) Interested(sender bool) bool {
	if sender {
		return c.SenderInterested
	}

	return c.ContactInterested
}

type Response string

const (
	None         Response = "none"
	Accepted     Response = "accepted"
	Declined     Response = "declined"
	FakeAccepted Response = "fake_accepted"
)

func (r Response) IsValid() bool {
	switch r {
	case None, Accepted, Declined, FakeAccepted:
		return true
	}

	return false
}
