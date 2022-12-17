package relate

func NewUser(db Database) *User {
	return &User{db: db}
}

type User struct {
	db Database
}

func (n User) MyGroups() memberOf {
	return memberOf(n)
}
