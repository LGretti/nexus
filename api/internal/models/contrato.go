package models

import "time"

type Contrato struct {
	ID               int64     `json:"id"`
	EmpresaID        int64     `json:"empresa_id"`
	TipoContrato     string    `json:"tipo_contrato"`
	HorasContratadas int       `json:"horas_contratadas"`
	DataInicio       time.Time `json:"data_inicio"`
	DataFim          time.Time `json:"data_fim"`
	Ativo            bool      `json:"ativo"`
}
