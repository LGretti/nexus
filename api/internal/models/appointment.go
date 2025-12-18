package models

import "time"

type Appointment struct {
	ID         int64 `json:"id" db:"id"`
	ContractID int64 `json:"contractId" db:"contract_id"`
	UserID     int64 `json:"userId" db:"user_id"`

	StartTime   time.Time  `json:"startTime" db:"start_time"`
	EndTime     *time.Time `json:"endTime" db:"end_time"`
	Description string     `json:"description" db:"description"`

	// Calculadas
	ContractTitle string    `json:"contractTitle,omitempty"` // Para mostrar "Ademicon" no grid
	UserName      string    `json:"userName,omitempty"`      // Para mostrar "Lucas"
	TotalHours    float64   `json:"totalHours"`              // Calculado (Fim - Início)
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

func (a *Appointment) GetID() int64 {
	return a.ID
}

func (a *Appointment) SetID(id int64) {
	a.ID = id
}
