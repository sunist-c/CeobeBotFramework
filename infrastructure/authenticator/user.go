package authenticator

type IUser interface {
	Key() string
}

type User string

func (u User) Key() string {
	return string(u)
}

func NewUser(key string) IUser {
	return User(key)
}
