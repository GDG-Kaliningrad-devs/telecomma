package user

type Provider interface {
	Get(id int) User
}

type provider struct {
	users map[int]User
}

func (p provider) Get(id int) User {
	u, ok := p.users[id]
	if ok {
		return u
	}

	return User{
		ID:       id,
		Name:     "(пользователь скрыл себя)",
		UserName: "?",
	}
}

func NewProvider(users []User) Provider {
	index := map[int]User{}

	for _, user := range users {
		index[user.ID] = user
	}

	return provider{users: index}
}
