package repository

import (
	"context"
	"database/sql"
	"fmt"

	"nexus/api/internal/models"
)

// UsuarioRepository define a interface para as operações com usuários.
// Usar uma interface é uma boa prática para permitir testes e "mocking".
type UsuarioRepository interface {
	Save(usuario models.Usuario) (models.Usuario, error)
	Get(id *int64) ([]models.Usuario, error)
	Update(usuario models.Usuario) (int64, error)
	Delete(id int64) (int64, error)
}

// postgresUsuarioRepository é a implementação da interface para o PostgreSQL.
type postgresUsuarioRepository struct {
	db *sql.DB
}

// NewUsuarioRepository cria uma nova instância do repositório de usuários.
func NewUsuarioRepository(db *sql.DB) UsuarioRepository {
	return &postgresUsuarioRepository{
		db: db,
	}
}

func (r *postgresUsuarioRepository) Save(usuario models.Usuario) (models.Usuario, error) {
	var exists bool
	err := r.db.QueryRowContext(context.Background(), "SELECT EXISTS(SELECT 1 FROM usuarios WHERE email = $1)", usuario.Email).Scan(&exists)
	if err != nil {
		return models.Usuario{}, fmt.Errorf("erro ao checar e-mail: %w", err)
	}
	if exists {
		return models.Usuario{}, fmt.Errorf("e-mail já cadastrado")
	}

	query := `INSERT INTO usuarios (nome, email, perfil) VALUES ($1, $2, $3) RETURNING id`
	err = r.db.QueryRowContext(context.Background(), query, usuario.Nome, usuario.Email, usuario.Perfil).Scan(&usuario.ID)
	if err != nil {
		return models.Usuario{}, err
	}
	return usuario, nil
}

func (r *postgresUsuarioRepository) Get(id *int64) ([]models.Usuario, error) {
	query := "SELECT id, nome, email, perfil FROM usuarios"
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

	var usuarios []models.Usuario
	for rows.Next() {
		var usuario models.Usuario
		if err := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.Perfil); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, usuario)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(usuarios) == 0 {
		return nil, fmt.Errorf("nenhum usuário encontrado")
	}
	return usuarios, nil
}

func (r *postgresUsuarioRepository) Update(usuario models.Usuario) (int64, error) {
	query := `UPDATE usuarios SET nome = $1, email = $2, perfil = $3 WHERE id = $4`
	res, err := r.db.ExecContext(context.Background(), query, usuario.Nome, usuario.Email, usuario.Perfil, usuario.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *postgresUsuarioRepository) Delete(id int64) (int64, error) {
	query := `DELETE FROM usuarios WHERE id = $1`
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
