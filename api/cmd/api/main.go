// O pacote main indica que este é o ponto de partida executável do nosso programa.
package main

import (
	"log"
	"net/http"

	"nexus/api/internal/handlers"
)

func main() {
	// Define a porta do servidor
	const port = ":8080"

	router := handlers.NewRouter()


	log.Printf("🚀 Servidor subindo na porta %s", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}