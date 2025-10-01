package repository

import (
	"context"
	"database/sql"

	"nexus/api/internal/models"
)

// EmpresaRepository define a interface para as operações com empresas.
type EmpresaRepository interface {
	Repository[*models.Empresa]
	SaveBatch(empresas []*models.Empresa) ([]*models.Empresa, error)
}

// postgresEmpresaRepository é a implementação da interface para o PostgreSQL.
type postgresEmpresaRepository struct {
	Repository[*models.Empresa]
	db *sql.DB
}

// NewEmpresaRepository cria uma nova instância do repositório de empresas.
func NewEmpresaRepository(db *sql.DB) EmpresaRepository {
	return &postgresEmpresaRepository{
		Repository: NewPostgresRepository[*models.Empresa](db, "empresas"),
		db:         db,
	}
}

// SaveBatch salva uma lista de empresas em uma única transação.
func (r *postgresEmpresaRepository) SaveBatch(empresas []*models.Empresa) ([]*models.Empresa, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(context.Background(), `INSERT INTO empresas (nome, cnpj, email_contato) VALUES ($1, $2, $3) RETURNING id`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var empresasSalvas []*models.Empresa
	for _, empresa := range empresas {
		var newID int64
		err := stmt.QueryRowContext(context.Background(), empresa.Nome, empresa.CNPJ, empresa.EmailContato).Scan(&newID)
		if err != nil {
			return nil, err
		}
		empresa.SetID(newID)
		empresasSalvas = append(empresasSalvas, empresa)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return empresasSalvas, nil
}
