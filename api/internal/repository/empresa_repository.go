package repository

import (
	"context"
	"database/sql"
	"fmt"

	"nexus/api/internal/models"
)

// EmpresaRepository define a interface para as operações com empresas.
// Usar uma interface é uma boa prática para permitir testes e "mocking".
type EmpresaRepository interface {
	Save(empresa models.Empresa) (models.Empresa, error)
	Get(id *int64) ([]models.Empresa, error)
	Update(empresa models.Empresa) (int64, error)
	Delete(id int64) (int64, error)
}

// postgresEmpresaRepository é a implementação da interface para o PostgreSQL.
type postgresEmpresaRepository struct {
	db *sql.DB
}

// NewEmpresaRepository cria uma nova instância do repositório de empresas.
func NewEmpresaRepository(db *sql.DB) EmpresaRepository {
	return &postgresEmpresaRepository{
		db: db,
	}
}

func (r *postgresEmpresaRepository) Save(empresa models.Empresa) (models.Empresa, error) {
	empresas := []models.Empresa{empresa}
	empresasSalvas, err := r.SaveBatch(empresas)
	if err != nil {
		return models.Empresa{}, err
	}
	return empresasSalvas[0], nil
}

func (r *postgresEmpresaRepository) SaveBatch(empresas []models.Empresa) ([]models.Empresa, error) {
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

	var empresasSalvas []models.Empresa

	for _, empresa := range empresas {
		var newID int64
		err := stmt.QueryRowContext(context.Background(), empresa.Nome, empresa.CNPJ, empresa.EmailContato).Scan(&newID)
		if err != nil {
			return nil, err
		}
		empresa.ID = int(newID)
		empresasSalvas = append(empresasSalvas, empresa)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return empresasSalvas, nil
}

func (r *postgresEmpresaRepository) Get(id *int64) ([]models.Empresa, error) {
	query := "SELECT id, nome, cnpj, email_contato FROM empresas"
	args := []interface{}{}
	if id != nil {
		query += " WHERE id = $1"
		args = append(args, *id)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var empresas []models.Empresa
	for rows.Next() {
		var emp models.Empresa
		if err := rows.Scan(&emp.ID, &emp.Nome, &emp.CNPJ, &emp.EmailContato); err != nil {
			return nil, err
		}
		empresas = append(empresas, emp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(empresas) == 0 {
		return nil, fmt.Errorf("nenhuma empresa encontrada")
	}
	return empresas, nil
}

func (r *postgresEmpresaRepository) Update(empresa models.Empresa) (int64, error) {
	query := `UPDATE empresas SET nome = $1, cnpj = $2, email_contato_faturamento = $3 WHERE id = $4`

	// ExecContext é usado para queries que não retornam linhas (como UPDATE, DELETE).
	res, err := r.db.ExecContext(context.Background(), query, empresa.Nome, empresa.CNPJ, empresa.EmailContato, empresa.ID)
	if err != nil {
		return 0, err
	}

	// RowsAffected retorna o número de linhas que foram afetadas pela query.
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	// Se nenhuma linha foi afetada, significa que a empresa com aquele ID não foi encontrada.
	return rowsAffected, nil
}

func (r *postgresEmpresaRepository) Delete(id int64) (int64, error) {
	query := `DELETE FROM empresas WHERE id = $1`

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
