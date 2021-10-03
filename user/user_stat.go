package user

type Stat struct {
	User
	Place            uint
	Contacts         []ContactStatus
	DeclinesCount    uint
	FakeAcceptsCount uint
}

type StatList []Stat

func (s StatList) IDs() []int {
	ids := make([]int, len(s))

	for i := range s {
		ids[i] = s[i].ID
	}

	return ids
}

func (s StatList) Len() int {
	return len(s)
}

func (s StatList) Less(i, j int) bool {
	return len(s[j].Contacts) < len(s[i].Contacts)
}

func (s StatList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
