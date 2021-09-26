package contact

import (
	"sort"
	"sync"
	"time"

	"gdg-kld.ru/telecomma/user"
	"gorm.io/gorm"
)

type statCalculator struct {
	m               sync.Mutex
	calculationTime time.Time
	result          []user.Stat
}

func (c *statCalculator) top(db *gorm.DB) (result []user.Stat, t time.Time, err error) {
	c.m.Lock()

	if time.Since(c.calculationTime) > 5*time.Minute {
		result, err = fetchStats(db)
		if err != nil {
			return result, c.calculationTime, err
		}

		c.result = result
		c.calculationTime = time.Now()
	}

	c.m.Unlock()

	return c.result, c.calculationTime, nil
}

//nolint:funlen,cyclop // refactor priority 1:split,optimize
func fetchStats(db *gorm.DB) ([]user.Stat, error) {
	var (
		contacts []user.Contact
		stats    = map[int]user.Stat{}
	)

	err := db.Find(&contacts).Error
	if err != nil {
		return nil, err
	}

	for _, contact := range contacts {
		senderStat, ok := stats[contact.SenderID]
		if !ok {
			senderStat = user.Stat{ID: contact.SenderID}
		}

		receiverStat, ok := stats[contact.ContactID]
		if !ok {
			receiverStat = user.Stat{ID: contact.ContactID}
		}

		switch contact.Response {
		case user.None:
		case user.Accepted:
			senderStat.ContactsCount++
			receiverStat.ContactsCount++

		case user.Declined:
			senderStat.DeclinesCount++

		case user.FakeAccepted:
			senderStat.FakeAcceptsCount++
		}

		stats[contact.SenderID] = senderStat
		stats[contact.ContactID] = receiverStat
	}

	statsSlice := make(user.StatList, len(stats))

	i := 0

	for _, stat := range stats {
		statsSlice[i] = stat
		i++
	}

	sort.Sort(statsSlice)

	if len(statsSlice) > 8 {
		statsSlice = statsSlice[0:8]
	}

	var users []user.User

	err = db.Find(&users, statsSlice.IDs()).Error
	if err != nil {
		return nil, err
	}

	statsSlice = statsSlice.WithNames(users)

	return statsSlice, nil
}
