package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB() (*sql.DB, error) {
	dsn := "postgres://nexususer:postgres@localhost:5432/nexusdb?sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Houve um erro com a conexão com o banco de dados:", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao pinar o banco: %w", err)
	}

	fmt.Println("Conexão com o banco de dados estabelecida")
	return db, nil
}
