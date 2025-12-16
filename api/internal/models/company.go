package models

type Company struct {
	ID           int64  `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	CNPJ         string `json:"cnpj" db:"cnpj"`
	ContactEmail string `json:"email" db:"contact_email"`
}

func (c *Company) GetID() int64 {
	return c.ID
}

func (c *Company) SetID(id int64) {
	c.ID = id
}
