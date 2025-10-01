package models

type Usuario struct {
	ID     int64  `json:"id"`
	Nome   string `json:"nome"`
	Email  string `json:"email"`
	Perfil string `json:"perfil"`
}

func (u *Usuario) GetID() int64 {
	return u.ID
}

func (u *Usuario) SetID(id int64) {
	u.ID = id
}
