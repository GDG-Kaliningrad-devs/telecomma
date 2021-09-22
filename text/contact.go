package text

import "fmt"

const (
	SearchResult             = `todo: search result`
	SearchNotFound           = `todo: no results`
	ContactRequestSend       = `todo: send`
	ContactRequestNotID      = `todo: not id`
	contactRequest           = `todo: contact request from user %s`
	ContactRequestApproveBtn = `Подтвердить`
	ContactRequestDeclineBtn = `Отклонить`
	ContactRequestErrSent    = `todo: request was already sent`
	ContactRequestSuccess    = `todo: contact created`
	ContactRequestDeclined   = `todo: declined warning`
)

func ContactRequest(user string) string {
	return fmt.Sprintf(contactRequest, user)
}
