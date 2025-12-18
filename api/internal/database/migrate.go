package database

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB) {
	// Cria instância do driver Postgres
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Não foi possível criar driver do banco: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", // <--- NOME DO DRIVER MUDOU
		driver,
	)
	if err != nil {
		log.Fatalf("Erro ao configurar migração: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Erro ao aplicar migrações: %v", err)
	}

	log.Println("✅ Migrations (Postgres) aplicadas com sucesso!")
}
