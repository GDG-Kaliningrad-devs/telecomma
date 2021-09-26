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
	LookAtTop = `Теперь можно посмотреть топ, найди там себя

/top`
)

func Top(list user.StatList, t time.Time) string {
	results := ""

	for _, stat := range list {
		results += fmt.Sprintf(
			"%d - %s, @%s, знакомств: %d, fakes: %d\n",
			stat.Place, stat.Name, stat.UserName, stat.ContactsCount, stat.FakeAcceptsCount,
		)
	}

	return fmt.Sprintf(top, results, t.Format("15:04"))
}
