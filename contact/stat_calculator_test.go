package contact //nolint:testpackage // internal test

import (
	"log"
	"testing"

	"gdg-kld.ru/telecomma/text"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//nolint:paralleltest // runs locally
func Test_statCalculator_top(t *testing.T) {
	t.Skip("local testing, comment to run")

	db, err := gorm.Open(sqlite.Open("../telecomma.db"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	c := statCalculator{}

	top, resultTime, err := c.top(db)
	assert.NoError(t, err)

	if err != nil {
		return
	}

	assert.Len(t, top, 8)

	if len(top) > 1 {
		assert.True(t, top[0].ContactsCount >= top[1].ContactsCount)
		assert.NotEqual(t, 0, top[0].ContactsCount)
	}

	for _, stat := range top {
		log.Println(stat)
	}

	log.Println(text.Top(top, resultTime))
}
