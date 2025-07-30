package saveuser

import "os/user"

type StoreUser interface {
	Save(user user.User) error
}
