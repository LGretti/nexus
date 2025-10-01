package repository

import (
	"context"
	"database/sql"
	"fmt"
	"nexus/api/internal/models"
)

// UsuarioRepository define a interface para as operações com usuários.
type UsuarioRepository interface {
	Repository[*models.Usuario]
	EmailExists(email string) (bool, error)
}

// postgresUsuarioRepository é a implementação da interface para o PostgreSQL.
type postgresUsuarioRepository struct {
	Repository[*models.Usuario]
	db *sql.DB
}

// NewUsuarioRepository cria uma nova instância do repositório de usuários.
func NewUsuarioRepository(db *sql.DB) UsuarioRepository {
	return &postgresUsuarioRepository{
		Repository: NewPostgresRepository[*models.Usuario](db, "usuarios"),
		db:         db,
	}
}

// EmailExists verifica se um e-mail já está cadastrado.
func (r *postgresUsuarioRepository) EmailExists(email string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE email = $1)", r.GetTableName())
	err := r.db.QueryRowContext(context.Background(), query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("erro ao checar e-mail: %w", err)
	}
	return exists, nil
}
