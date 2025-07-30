package userinfo

import "os/user"

type UserProvider interface {
	User(int) (user.User, error)
}
