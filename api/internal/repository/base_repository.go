package repository

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"nexus/internal/models"
)

// Repository é uma interface para operações de banco de dados genéricas.
// O tipo T deve ser um ponteiro para uma struct que implementa models.Model.
type Repository[T models.Model] interface {
	Save(model T) (T, error)
	Get(id *int64) ([]T, error)
	Update(model T) (int64, error)
	Delete(id int64) (int64, error)
	GetTableName() string
}

// postgresRepository é a implementação da interface Repository para o PostgreSQL.
type postgresRepository[T models.Model] struct {
	db        *sql.DB
	tableName string
}

// NewPostgresRepository cria uma nova instância de postgresRepository.
func NewPostgresRepository[T models.Model](db *sql.DB, tableName string) Repository[T] {
	return &postgresRepository[T]{
		db:        db,
		tableName: tableName,
	}
}

func (r *postgresRepository[T]) GetTableName() string {
	return r.tableName
}

// Save insere um novo modelo no banco de dados.
func (r *postgresRepository[T]) Save(model T) (T, error) {
	val := reflect.ValueOf(model).Elem()
	typ := val.Type()

	var cols []string
	var values []interface{}
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == "ID" {
			continue
		}
		// Usamos o nome da tag db como nome da coluna
		dbTag := strings.Split(field.Tag.Get("db"), ",")[0]
		if dbTag != "" {
			cols = append(cols, dbTag)
			values = append(values, val.Field(i).Interface())
		}
	}

	colNames := strings.Join(cols, ", ")
	placeholders := ""
	for i := range cols {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", r.tableName, colNames, placeholders)

	var id int64
	err := r.db.QueryRowContext(context.Background(), query, values...).Scan(&id)
	if err != nil {
		return model, fmt.Errorf("erro ao inserir no banco de dados: %w", err)
	}

	model.SetID(id)
	return model, nil
}

// Get recupera um ou todos os modelos do banco de dados.
func (r *postgresRepository[T]) Get(id *int64) ([]T, error) {
	var t T
	typ := reflect.TypeOf(t).Elem()

	var cols []string
	for i := 0; i < typ.NumField(); i++ {
		dbTag := strings.Split(typ.Field(i).Tag.Get("db"), ",")[0]
		if dbTag != "" {
			cols = append(cols, dbTag)
		}
	}
	colNames := strings.Join(cols, ", ")

	query := fmt.Sprintf("SELECT %s FROM %s", colNames, r.tableName)
	args := []interface{}{}
	if id != nil {
		query += " WHERE id = $1"
		args = append(args, *id)
	}

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar o banco de dados: %w", err)
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		newElemPtr := reflect.New(typ)
		result := newElemPtr.Interface().(T)

		resultValue := newElemPtr.Elem()
		scanArgs := make([]interface{}, resultValue.NumField())
		for i := 0; i < resultValue.NumField(); i++ {
			scanArgs[i] = resultValue.Field(i).Addr().Interface()
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("erro ao escanear linha: %w", err)
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro nas linhas: %w", err)
	}

	return results, nil
}

// Update atualiza um modelo existente no banco de dados.
func (r *postgresRepository[T]) Update(model T) (int64, error) {
	val := reflect.ValueOf(model).Elem()
	typ := val.Type()

	var setClauses []string
	var values []interface{}
	argCount := 1
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == "ID" {
			continue
		}
		dbTag := strings.Split(field.Tag.Get("db"), ",")[0]
		if dbTag != "" {
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbTag, argCount))
			values = append(values, val.Field(i).Interface())
			argCount++
		}
	}

	values = append(values, model.GetID())

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", r.tableName, strings.Join(setClauses, ", "), argCount)
	res, err := r.db.ExecContext(context.Background(), query, values...)
	if err != nil {
		return 0, fmt.Errorf("erro ao atualizar no banco de dados: %w", err)
	}
	return res.RowsAffected()
}

// Delete remove um modelo do banco de dados.
func (r *postgresRepository[T]) Delete(id int64) (int64, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", r.tableName)
	res, err := r.db.ExecContext(context.Background(), query, id)
	if err != nil {
		return 0, fmt.Errorf("erro ao deletar no banco de dados: %w", err)
	}
	return res.RowsAffected()
}
