package user

import (
	"strings"

	"gopkg.in/tucnak/telebot.v2"
)

func Name(user *telebot.User) string {
	return strings.Join([]string{user.FirstName, user.LastName}, " ")
}
