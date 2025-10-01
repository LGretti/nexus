package models

type Empresa struct {
	ID           int64  `json:"id"`
	Nome         string `json:"nome"`
	CNPJ         string `json:"cnpj"`
	EmailContato string `json:"email_contato"`
}

func (e *Empresa) GetID() int64 {
	return e.ID
}

func (e *Empresa) SetID(id int64) {
	e.ID = id
}
