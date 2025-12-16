package repository

import (
	"context"
	"database/sql"
	"fmt"
	"nexus/api/internal/models"
)

// UserRepository define a interface para as operações com usuários.
type UserRepository interface {
	Repository[*models.User]
	EmailExists(email string) (bool, error)
}

// postgresUserRepository é a implementação da interface para o PostgreSQL.
type postgresUserRepository struct {
	Repository[*models.User]
	db *sql.DB
}

// NewUserRepository cria uma nova instância do repositório de usuários.
func NewUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{
		Repository: NewPostgresRepository[*models.User](db, "users"),
		db:         db,
	}
}

// EmailExists verifica se um e-mail já está cadastrado.
func (r *postgresUserRepository) EmailExists(email string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE email = $1)", r.GetTableName())
	err := r.db.QueryRowContext(context.Background(), query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("erro ao checar e-mail: %w", err)
	}
	return exists, nil
}
