package domain

// PgUser TODO: remove pointers. Save omitempty fields like ""
type PgUser struct {
	ID        string  `json:"id"`
	TgUserId  int64   `json:"tg_uid"`
	UserName  *string `json:"login,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}
