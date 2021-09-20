package user

type User struct {
	ID       int // telegram user id
	Name     string
	Bio      string  // short self description
	Contacts []*User `gorm:"many2many:user_contacts"`
}

func NewUser(id int, name string) (User, error) {
	name, err := validateName(name)
	if err != nil {
		return User{}, err
	}

	return User{ID: id, Name: name}, nil
}
