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
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close() // Fecha conexão se o main for de f

	// 2. Auto-Migration
	database.RunMigrations(db)

	// 3. Repositórios
	contractRepo := repository.NewContractRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	userRepo := repository.NewUserRepository(db)
	appointmentRepo := repository.NewAppointmentRepository(db)

	// 4. Handlers
	companyHandler := handlers.NewCompanyHandler(companyRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	contractHandler := handlers.NewContractHandler(contractRepo)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentRepo)

	// 5. Roteador
	router := handlers.NewRouter(companyHandler, userHandler, contractHandler, appointmentHandler)

	const port = ":8080"
	log.Printf("Servidor subindo na porta %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
