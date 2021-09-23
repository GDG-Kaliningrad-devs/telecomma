package text

import "fmt"

const (
	ActionDone                  = `Выполнено`
	SearchResult                = `Вы можете создать контакт с этими пользователями`
	SearchErrMinSymbols         = `Минимум 4 символа`
	SearchErrNoResults          = `Никого с таким именем не найдено. Возможно, бота не запускали или не передали своё имя`
	SearchErrToManyResults      = `Найдено слишком много, напишите ещё пару символов`
	ContactRequestSend          = `Запрос отправлен`
	ContactRequestWrongData     = `Error: wrong contact data. Pls contact org`
	contactRequest              = `%s, %s хочет сохранить контакт`
	ContactRequestApproveBtn    = `Подтвердить`
	ContactRequestDeclineBtn    = `Отклонить`
	ContactRequestFakeAcceptBtn = `Притвориться, что подтвердил`
	ContactRequestErrSent       = `Вы уже отправляли запрос`
	ContactRequestErrIgnored    = `Собеседник уже отправлял запрос вам`
	ContactResponseErrSent      = `Вы уже ответили на запрос`
	contactRequestSuccess       = `#contact
%s @%s и %s @%s теперь дружбаны`
	contactRequestDeclined = `%s отклонил запрос. Общайтесь по-настоящему и спрашивайте разрешения`
	ContactSetImportantBtn = `Хочу пообщаться после конфы`
)

func ContactRequest(user, username string) string {
	return fmt.Sprintf(contactRequest, user, username)
}

func ContactRequestDeclined(user string) string {
	return fmt.Sprintf(contactRequestDeclined, user)
}

func ContactRequestSuccess(name1, username1, name2, username2 string) string {
	return fmt.Sprintf(contactRequestSuccess, name1, username1, name2, username2)
}
