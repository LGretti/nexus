package models

import "time"

type Contract struct {
	ID           int64     `json:"id" db:"id"`
	CompanyId    int64     `json:"companyId" db:"company_id"`
	CompanyName  string    `json:"companyName,omitempty" db:"company_name,omitempty"`
	Title        string    `json:"title" db:"title"`
	ContractType string    `json:"contractType" db:"contract_type"`
	TotalHours   int       `json:"totalHours" db:"total_hours"`
	StartDate    time.Time `json:"startDate" db:"start_date"`
	EndDate      time.Time `json:"endDate" db:"end_date"`
	IsActive     bool      `json:"isActive" db:"is_active"`
}

func (c *Contract) GetID() int64 {
	return c.ID
}

func (c *Contract) SetID(id int64) {
	c.ID = id
}
