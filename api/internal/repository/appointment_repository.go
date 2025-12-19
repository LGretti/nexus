package repository

import (
	"context"
	"database/sql"
	"nexus/internal/models"
)

type AppointmentRepository interface {
	Repository[*models.Appointment]
	GetAllWithContract() ([]*models.Appointment, error)
	GetByContractID(contractID int64) ([]*models.Appointment, error)
	GetByUserID(userID int64) ([]*models.Appointment, error)
}

type postgresAppointmentRepository struct {
	Repository[*models.Appointment]
	db *sql.DB
}

func NewAppointmentRepository(db *sql.DB) AppointmentRepository {
	return &postgresAppointmentRepository{
		Repository: NewPostgresRepository[*models.Appointment](db, "appointments"),
		db:         db,
	}
}

func scanAppointments(rows *sql.Rows) ([]*models.Appointment, error) {
	var appointments []*models.Appointment
	for rows.Next() {
		var a models.Appointment
		if err := rows.Scan(
			&a.ID, &a.StartTime, &a.EndTime, &a.Description, &a.TotalHours, &a.DurationSeconds,
			&a.ContractTitle, &a.UserName,
		); err != nil {
			return nil, err
		}
		appointments = append(appointments, &a)
	}
	return appointments, nil
}

func (r *postgresAppointmentRepository) GetAllWithContract() ([]*models.Appointment, error) {
	query := `SELECT a.id, a.start_time, a.end_time, a.description, -- Adicionei description
	                 EXTRACT(EPOCH FROM (COALESCE(a.end_time, CURRENT_TIMESTAMP) - a.start_time)) / 3600 as total_hours,
									 EXTRACT(EPOCH FROM (COALESCE(a.end_time, CURRENT_TIMESTAMP) - a.start_time))::bigint as duration_seconds,
                     c.title, u.name
              FROM appointments a
              JOIN contracts c ON a.contract_id = c.id
              JOIN users u ON a.user_id = u.id
              ORDER BY a.created_at DESC`

	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAppointments(rows)
}

func (r *postgresAppointmentRepository) GetByContractID(contractID int64) ([]*models.Appointment, error) {
	query := `SELECT a.id, a.start_time, a.end_time, a.description,
	                 EXTRACT(EPOCH FROM (COALESCE(a.end_time, CURRENT_TIMESTAMP) - a.start_time)) / 3600 as total_hours,
									 EXTRACT(EPOCH FROM (COALESCE(a.end_time, CURRENT_TIMESTAMP) - a.start_time))::bigint as duration_seconds,
                     c.title, u.name
              FROM appointments a
              JOIN contracts c ON a.contract_id = c.id
              JOIN users u ON a.user_id = u.id
              WHERE a.contract_id = $1
              ORDER BY a.created_at DESC`

	rows, err := r.db.QueryContext(context.Background(), query, contractID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAppointments(rows)
}

func (r *postgresAppointmentRepository) GetByUserID(userID int64) ([]*models.Appointment, error) {
	query := `
		SELECT id, contract_id, user_id, description, start_time, end_time, created_at
		FROM appointments
		WHERE user_id = $1
		ORDER BY start_time DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var appointments []*models.Appointment

	for rows.Next() {
		var appt models.Appointment
		// Lembre-se: appt.EndTime é um ponteiro (*time.Time) para aceitar NULL do banco
		if err := rows.Scan(
			&appt.ID,
			&appt.ContractID,
			&appt.UserID,
			&appt.Description,
			&appt.StartTime,
			&appt.EndTime,
			&appt.CreatedAt,
		); err != nil {
			return nil, err
		}
		appointments = append(appointments, &appt)
	}
	return appointments, nil
}
