package text

import (
	"fmt"
	"time"

	"gdg-kld.ru/telecomma/user"
)

const (
	top = `Вот самые общительные ребята конференции:

%s

Результат на: %s

Успейте попасть в топ до конца конфы`
	LookAtTop = `Пофиксил, теперь точно можно узнать своё место

(👍≖‿‿≖)👍 /top`
	stats = `
Вот статистика функции "Хочу пообщаться после конфы"
Вы совпали c: %d
%s
С вами хотят пообщаться: %d
%s`
	finalMessage = `Привет! Моё (бота) посмертное уведомление.

Через неделю бот перестанет работать. Пока вы ещё можете найти и связаться с кем-то

Вы можете посмотреть контакты с конференции нажав на тег #contact
%s
Вы можете помочь порефакторить и запилить полезных функций для следующих митапов, хакатонов и фестов тут
github.com/GDG-Kaliningrad-devs/telecomma
Просто создайте issue "what to do"

ʕ•́ᴥ•̀ʔっ♡`
)

func Top(list user.StatList, t time.Time) string {
	results := ""

	for _, stat := range list {
		results += fmt.Sprintf(
			"%d - %s, @%s, знакомств: %d, fakes: %d\n",
			stat.Place, stat.Name, stat.UserName, len(stat.Contacts), stat.FakeAcceptsCount,
		)
	}

	return fmt.Sprintf(top, results, t.Format("15:04"))
}

func Final(contacts []user.ContactStatus) string {
	var (
		matchesCount    uint
		matches         string
		interestedCount uint
		interested      string
	)

	for _, contact := range contacts {
		switch contact.MatchStatus {
		case user.Match:
			matches += contact.User.String() + "\n"
			matchesCount++

		case user.InterestedInMe:
			interested += contact.User.String() + "\n"
			interestedCount++

		case user.Nothing:
		}
	}

	if len(matches) == 0 && len(interested) == 0 {
		return fmt.Sprintf(finalMessage, "")
	}

	statStr := fmt.Sprintf(
		stats,
		matchesCount, matches,
		interestedCount, interested,
	)

	return fmt.Sprintf(finalMessage, statStr)
}
