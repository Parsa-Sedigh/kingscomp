package entity

import "time"

type Account struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	UserName    string    `json:"user_name"`
	DisplayName string    `json:"display_name"`
	JoinedAt    time.Time `json:"joined_at"`
	State       string    `json:"state"`
}

func (a Account) EntityID() ID {
	return NewID("account", a.ID)
}
