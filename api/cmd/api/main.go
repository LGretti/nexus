package main

import (
	"log"
	"net/http"

	"nexus/api/internal/database"
	"nexus/api/internal/handlers"
)

func main() {
	const port = ":8080"

	database.ConnectDB()

	router := handlers.NewRouter()

	log.Printf("Iniciando servidor na porta %s", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
