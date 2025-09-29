package repository

import (
	"context"
	"database/sql"
	"fmt"

	"nexus/api/internal/models"
)

// ContratoRepository define a interface para as operações com contratos.
// Usar uma interface é uma boa prática para permitir testes e "mocking".
type ContratoRepository interface {
	Save(contrato models.Contrato) (models.Contrato, error)
	Get(id *int64) ([]models.Contrato, error)
	Update(contrato models.Contrato) (int64, error)
	Delete(id int64) (int64, error)
	GetPorEmpresaID(empresaID int64) ([]models.Contrato, error)
}

// postgresContratoRepository é a implementação da interface para o PostgreSQL.
type postgresContratoRepository struct {
	db *sql.DB
}

// NewContratoRepository cria uma nova instância do repositório de contratos.
func NewContratoRepository(db *sql.DB) ContratoRepository {
	return &postgresContratoRepository{
		db: db,
	}
}

func (r *postgresContratoRepository) Save(contrato models.Contrato) (models.Contrato, error) {
	query := `INSERT INTO contratos (empresa_id, tipo_contrato, horas_contratadas, data_inicio, data_fim) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	if contrato.DataFim.Before(contrato.DataInicio) {
		return models.Contrato{}, fmt.Errorf("data de fim não pode ser anterior à data de início")
	}

	err := r.db.QueryRowContext(context.Background(), query, contrato.EmpresaID, contrato.TipoContrato, contrato.HorasContratadas, contrato.DataInicio, contrato.DataFim).Scan(&contrato.ID)
	if err != nil {
		return models.Contrato{}, err
	}
	return contrato, nil
}

func (r *postgresContratoRepository) Get(id *int64) ([]models.Contrato, error) {
	query := "SELECT id, empresa_id, tipo_contrato, horas_contratadas, data_inicio, data_fim FROM contratos"
	args := []interface{}{}
	if id != nil {
		query += " WHERE id = $1"
		args = append(args, *id)
	}

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contratos []models.Contrato
	for rows.Next() {
		var contrato models.Contrato
		if err := rows.Scan(&contrato.ID, &contrato.EmpresaID, &contrato.TipoContrato, &contrato.HorasContratadas, &contrato.DataInicio, &contrato.DataFim); err != nil {
			return nil, err
		}
		contratos = append(contratos, contrato)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(contratos) == 0 {
		return nil, fmt.Errorf("nenhum contrato encontrado")
	}
	return contratos, nil
}

func (r *postgresContratoRepository) Update(contrato models.Contrato) (int64, error) {
	query := `UPDATE contratos SET empresa_id = $1, tipo_contrato = $2, horas_contratadas = $3, data_inicio = $4, data_fim = $5 WHERE id = $6`

	if contrato.DataFim.Before(contrato.DataInicio) {
		return 0, fmt.Errorf("data de fim não pode ser anterior à data de início")
	}

	err := r.db.QueryRowContext(context.Background(), query, contrato.EmpresaID, contrato.TipoContrato, contrato.HorasContratadas, contrato.DataInicio, contrato.DataFim, contrato.ID).Scan(&contrato.ID)
	if err != nil {
		return 0, err
	}
	return contrato.ID, nil
}

// TODO: Contrato não deve ser deletado, mas gerar um campo "ativo" falso
func (r *postgresContratoRepository) Delete(id int64) (int64, error) {
	query := `UPDATE contratos SET ativo = false WHERE id = $1`
	res, err := r.db.ExecContext(context.Background(), query, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (r *postgresContratoRepository) GetPorEmpresaID(empresaID int64) ([]models.Contrato, error) {
	query := "SELECT id, empresa_id, tipo_contrato, horas_contratadas, data_inicio, data_fim FROM contratos WHERE empresa_id = $1"
	rows, err := r.db.QueryContext(context.Background(), query, empresaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contratos []models.Contrato
	for rows.Next() {
		var contrato models.Contrato
		if err := rows.Scan(&contrato.ID, &contrato.EmpresaID, &contrato.TipoContrato, &contrato.HorasContratadas, &contrato.DataInicio, &contrato.DataFim); err != nil {
			return nil, err
		}
		contratos = append(contratos, contrato)
	}
	return contratos, rows.Err()
}
