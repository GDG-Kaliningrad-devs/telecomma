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

func (c *statCalculator) ids(db *gorm.DB) ([]int, error) {
	err := c.assertCalculation(db)
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(c.result))

	for i := range c.result {
		ids[i] = c.result[i].ID
	}

	return ids, nil
}

func (c *statCalculator) top(db *gorm.DB, currentUserID int) (
	result []user.Stat,
	t time.Time,
	err error,
) {
	//
	err = c.assertCalculation(db)
	if err != nil {
		return nil, time.Time{}, err
	}

	var personalTop []user.Stat //nolint:prealloc

	for i, stat := range c.result {
		if i > 7 && stat.ID != currentUserID {
			continue
		}

		if stat.ID == currentUserID {
			stat.Name = "(вы) " + stat.Name
		}

		stat.Place = uint(i) + 1

		personalTop = append(personalTop, stat)
	}

	return personalTop, c.calculationTime, nil
}

func (c *statCalculator) assertCalculation(db *gorm.DB) error {
	c.m.Lock()

	if time.Since(c.calculationTime) > 5*time.Minute {
		result, err := fetchStats(db)
		if err != nil {
			return err
		}

		c.result = result
		c.calculationTime = time.Now()
	}

	c.m.Unlock()

	return nil
}

//nolint:cyclop // refactor priority 1 split,optimize
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

	var users []user.User

	err = db.Find(&users, statsSlice.IDs()).Error
	if err != nil {
		return nil, err
	}

	statsSlice = statsSlice.WithNames(users)

	return statsSlice, nil
}
