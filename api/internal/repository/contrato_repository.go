package repository

import (
	"context"
	"database/sql"
	"nexus/api/internal/models"
)

// ContratoRepository define a interface para as operações com contratos.
type ContratoRepository interface {
	Repository[*models.Contrato]
	GetPorEmpresaID(empresaID int64) ([]*models.Contrato, error)
	Delete(id int64) (int64, error)
}

// postgresContratoRepository é a implementação da interface para o PostgreSQL.
type postgresContratoRepository struct {
	Repository[*models.Contrato]
	db *sql.DB
}

// NewContratoRepository cria uma nova instância do repositório de contratos.
func NewContratoRepository(db *sql.DB) ContratoRepository {
	return &postgresContratoRepository{
		Repository: NewPostgresRepository[*models.Contrato](db, "contratos"),
		db:         db,
	}
}

func (r *postgresContratoRepository) GetPorEmpresaID(empresaID int64) ([]*models.Contrato, error) {
	query := "SELECT id, empresa_id, tipo_contrato, horas_contratadas, data_inicio, data_fim, ativo FROM contratos WHERE empresa_id = $1"
	rows, err := r.db.QueryContext(context.Background(), query, empresaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contratos []*models.Contrato
	for rows.Next() {
		var contrato models.Contrato
		if err := rows.Scan(&contrato.ID, &contrato.EmpresaID, &contrato.TipoContrato, &contrato.HorasContratadas, &contrato.DataInicio, &contrato.DataFim, &contrato.Ativo); err != nil {
			return nil, err
		}
		contratos = append(contratos, &contrato)
	}
	return contratos, rows.Err()
}

// Delete customizado para Contrato: desativa em vez de deletar.
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
