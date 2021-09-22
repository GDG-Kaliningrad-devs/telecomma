package user

import (
	"time"

	"gopkg.in/tucnak/telebot.v2"
)

type User struct {
	ID   int // telegram user id
	Name string
}

func NewUser(user *telebot.User) (User, error) {
	name, err := validateName(Name(user))
	if err != nil {
		return User{}, err
	}

	return User{ID: user.ID, Name: name}, nil
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

func (c Contact) Respond(accepted bool) Contact {
	if accepted {
		c.Response = Accepted
	} else {
		c.Response = Declined
	}

	now := time.Now()

	c.ResponseTime = &now

	return c
}

type Response string

const (
	None     Response = "none"
	Accepted Response = "accepted"
	Declined Response = "declined"
)

func (r Response) IsValid() bool {
	switch r {
	case None, Accepted, Declined:
		return true
	}

	return false
}
