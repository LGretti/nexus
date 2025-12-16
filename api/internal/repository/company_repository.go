package repository

import (
	"context"
	"database/sql"

	"nexus/api/internal/models"
)

// CompanyRepository define a interface para as operações com empresas.
type CompanyRepository interface {
	Repository[*models.Company]
	SaveBatch(companies []*models.Company) ([]*models.Company, error)
}

// postgresCompanyRepository é a implementação da interface para o PostgreSQL.
type postgresCompanyRepository struct {
	Repository[*models.Company]
	db *sql.DB
}

// NewCompanyRepository cria uma nova instância do repositório de empresas.
func NewCompanyRepository(db *sql.DB) CompanyRepository {
	return &postgresCompanyRepository{
		Repository: NewPostgresRepository[*models.Company](db, "companies"),
		db:         db,
	}
}

// SaveBatch salva uma lista de empresas em uma única transação.
func (r *postgresCompanyRepository) SaveBatch(companies []*models.Company) ([]*models.Company, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(context.Background(), `INSERT INTO companies (name
	                                                                            ,cnpj
																																					    ,contact_email)
																												VALUES ($1
																												       ,$2
																															 ,$3)
																												RETURNING id`,
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var companiesSaved []*models.Company
	for _, company := range companies {
		var newID int64
		err := stmt.QueryRowContext(context.Background(), company.Name, company.CNPJ, company.ContactEmail).Scan(&newID)
		if err != nil {
			return nil, err
		}
		company.SetID(newID)
		companiesSaved = append(companiesSaved, company)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return companiesSaved, nil
}
