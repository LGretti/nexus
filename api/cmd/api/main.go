package main

import (
	"log"
	"net/http"

	"nexus/api/internal/database"
	"nexus/api/internal/handlers"
	"nexus/api/internal/repository"
)

func main() {
	// 1. Conecta ao banco de dados
	database.ConnectDB()
	db := database.DB
	if db == nil {
		log.Fatalf("Não foi possível conectar ao banco de dados: DB is nil")
	}

	// 2. Cria as instâncias dos repositórios
	contratoRepo := repository.NewContratoRepository(db)
	empresaRepo := repository.NewEmpresaRepository(db)
	usuarioRepo := repository.NewUsuarioRepository(db)

	// 3. Cria as instâncias dos Handlers
	empresaHandler := handlers.NewEmpresaHandler(empresaRepo)
	usuarioHandler := handlers.NewUsuarioHandler(usuarioRepo)
	contratoHandler := handlers.NewContratoHandler(contratoRepo)

	// 4. Configura o roteador
	router := handlers.NewRouter(empresaHandler, usuarioHandler, contratoHandler)

	const port = ":8080"
	log.Printf("Servidor subindo na porta %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
