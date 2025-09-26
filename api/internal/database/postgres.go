package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"nexus/api/internal/models"
)

var DB *sql.DB

func ConnectDB() {
	dsn := "postgres://nexususer:postgres@localhost:5432/nexusdb?sslmode=disable"

	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Houve um erro com a conexão com o banco de dados:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Houve um erro ao tentar pingar o banco de dados:", err)
	}

	fmt.Println("Conexão com o banco de dados estabelecida com êxito")
}

func CreateEmpresa(empresa models.Empresa) (models.Empresa, error) {
	empresas := []models.Empresa{empresa}
	empresasSalvas, err := CreateEmpresasBatch(empresas)
	if err != nil {
		return models.Empresa{}, err
	}
	return empresasSalvas[0], nil

}

func CreateEmpresasBatch(empresas []models.Empresa) ([]models.Empresa, error) {
	tx, err := DB.Begin()
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

func GetEmpresas() ([]models.Empresa, error) {
	rows, err := DB.Query("SELECT id, nome, cnpj, email_contato FROM empresas")
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
	if len(empresas) == 0 {
		return nil, fmt.Errorf("nenhuma empresa encontrada")
	}
	return empresas, nil
}
