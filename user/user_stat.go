package user

type Stat struct {
	ID               int
	Name             string
	UserName         string
	ContactsCount    uint
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

func (s StatList) WithNames(users []User) StatList {
	for i := range s {
		stat := s[i]

		stat.Name = "(имя скрыто)"
		stat.UserName = "?"

		for _, u := range users {
			if u.ID != stat.ID {
				continue
			}

			stat.Name = u.Name
			stat.UserName = u.UserName
		}

		s[i] = stat
	}

	return s
}

func (s StatList) Len() int {
	return len(s)
}

func (s StatList) Less(i, j int) bool {
	return s[j].ContactsCount < s[i].ContactsCount
}

func (s StatList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
