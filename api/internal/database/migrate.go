package database

import (
	"database/sql"
	_ "embed"
	"log"
)

//go:embed schema.sql
var schemaSQL string

func RunMigrations(db *sql.DB) {
	log.Println("Validando estrutura do BD...")

	_, err := db.Exec(schemaSQL)
	if err != nil {
		log.Fatalf("Falha crítica ao criar tabelas: %v", err)
	}

	log.Println("Estrutura do BD atualizada")
}
