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

type Contact struct {
	SenderID     int `gorm:"primaryKey"`
	ContactID    int `gorm:"primaryKey"`
	Response     Response
	RequestTime  time.Time
	ResponseTime *time.Time
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
