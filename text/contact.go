package text

import "fmt"

const (
	ActionDone               = `Выполнено`
	SearchResult             = `todo: search result`
	SearchErrMinSymbols      = `Минимум 4 символа`
	SearchErrNoResults       = `Никого с таким именем не найдено. Возможно, бота не запускали или не передали своё имя`
	SearchErrToManyResults   = `Найдено слишком много, напишите ещё пару символов`
	ContactRequestSend       = `todo: send`
	ContactRequestNotID      = `todo: not id`
	contactRequest           = `todo: contact request from user %s`
	ContactRequestApproveBtn = `Подтвердить`
	ContactRequestDeclineBtn = `Отклонить`
	ContactRequestErrSent    = `todo: request was already sent`
	ContactRequestErrIgnored = `todo: request was already received`
	ContactResponseErrSent   = `todo: response was already sent`
	ContactRequestSuccess    = `todo: contact created`
	ContactRequestDeclined   = `todo: declined warning`
)

func ContactRequest(user string) string {
	return fmt.Sprintf(contactRequest, user)
}
