package models

type Empresa struct {
	ID           int    `json:"id"`
	Nome         string `json:"nome"`
	CNPJ         string `json:"cnpj"`
	EmailContato string `json:"email_contato"`
}
