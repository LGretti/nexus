package models

type User struct {
	ID    int64  `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Role  string `json:"role" db:"role"`
}

func (u *User) GetID() int64 {
	return u.ID
}

func (u *User) SetID(id int64) {
	u.ID = id
}
