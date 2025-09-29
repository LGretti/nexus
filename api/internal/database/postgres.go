package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
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
