package repository

import (
	"context"
	"database/sql"
	"fmt"
	"nexus/api/internal/models"
)

// ContractRepository define a interface para as operações com contratos.
type ContractRepository interface {
	Repository[*models.Contract]
	GetByCompanyID(companyID int64) ([]*models.Contract, error)
	GetAllWithCompany() ([]*models.Contract, error)
	Delete(id int64) (int64, error)
}

// postgresContractRepository é a implementação da interface para o PostgreSQL.
type postgresContractRepository struct {
	Repository[*models.Contract]
	db *sql.DB
}

// NewContractRepository cria uma nova instância do repositório de contratos.
func NewContractRepository(db *sql.DB) ContractRepository {
	return &postgresContractRepository{
		Repository: NewPostgresRepository[*models.Contract](db, "contracts"),
		db:         db,
	}
}

func (r *postgresContractRepository) GetAllWithCompany() ([]*models.Contract, error) {
	// AQUI ESTÁ O PULO DO GATO: O JOIN COM COMPANIES
	query := `
		SELECT contracts.id
		      ,contracts.company_id
					,companies.name
					,contracts.contract_type
					,contracts.total_hours
					,contracts.start_date
					,contracts.end_date
					,contracts.is_active
		FROM contracts
		     INNER JOIN companies
		     ON contracts.company_id = companies.id
		ORDER BY contracts.created_at DESC
	`

	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contracts []*models.Contract
	for rows.Next() {
		var c models.Contract
		if err := rows.Scan(
			&c.ID, &c.CompanyId, &c.CompanyName,
			&c.ContractType, &c.TotalHours, &c.StartDate, &c.EndDate, &c.IsActive,
		); err != nil {
			return nil, err
		}
		// Regra de negócio: Montar o Título para o Front
		c.Title = fmt.Sprintf("%s - %s", c.CompanyName, c.ContractType)

		contracts = append(contracts, &c)
	}
	return contracts, nil
}

func (r *postgresContractRepository) GetByCompanyID(companyID int64) ([]*models.Contract, error) {
	query := "SELECT id, company_id, contract_type, total_hours, start_date, end_date, is_active FROM contracts WHERE company_id = $1"
	rows, err := r.db.QueryContext(context.Background(), query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contracts []*models.Contract
	for rows.Next() {
		var contract models.Contract
		if err := rows.Scan(&contract.ID, &contract.CompanyId, &contract.ContractType, &contract.TotalHours, &contract.StartDate, &contract.EndDate, &contract.IsActive); err != nil {
			return nil, err
		}
		contracts = append(contracts, &contract)
	}
	return contracts, rows.Err()
}

// Delete customizado para Contrato: desativa em vez de deletar.
func (r *postgresContractRepository) Delete(id int64) (int64, error) {
	query := `UPDATE contracts SET is_active = false WHERE id = $1`
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
