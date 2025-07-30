package posgresql

type User struct {
	ID        int64  `json:"id"`
	UserName  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
}
