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

	top, resultTime, err := c.top(db, 591615813)
	assert.NoError(t, err)

	if err != nil {
		return
	}

	assert.GreaterOrEqual(t, len(top), 8)
	assert.LessOrEqual(t, len(top), 9)

	if len(top) > 1 {
		assert.True(t, len(top[0].Contacts) >= len(top[1].Contacts))
		assert.NotEqual(t, 0, len(top[0].Contacts))
	}

	for _, stat := range top {
		log.Println(stat)
	}

	log.Println(text.Top(top, resultTime))

	log.Println(text.Final(top[0].Contacts))
}
