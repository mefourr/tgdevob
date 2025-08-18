package posgresql

type User struct {
	ID        string `json:"id"`
	TgUserId  int    `json:"tg_uid"`
	UserName  string `json:"login,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}
