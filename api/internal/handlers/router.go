package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	companyHandler *CompanyHandler,
	userHandler *UserHandler,
	contractHandler *ContractHandler,
	appointmentHandler *AppointmentHandler, // Adicionado o novo handler
) http.Handler {

	r := chi.NewRouter()

	// Middlewares Globais (Logs, Recover, CORS)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Nexus API is running 🚀"))
	})

	// --- 1. ROTAS DE EMPRESAS (COMPANIES) ---
	r.Route("/api/companies", func(r chi.Router) {
		r.Post("/", companyHandler.CreateHandler)       // Criar empresa
		r.Get("/", companyHandler.GetAllHandler)        // Listar empresas
		r.Get("/{id}", companyHandler.GetByIDHandler)   // Detalhe da empresa
		r.Put("/{id}", companyHandler.UpdateHandler)    // Atualizar
		r.Delete("/{id}", companyHandler.DeleteHandler) // Deletar
	})

	// --- 2. ROTAS DE USUÁRIOS (USERS) ---
	r.Route("/api/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateHandler)
		r.Get("/", userHandler.GetAllHandler)
		r.Get("/{id}", userHandler.GetByIDHandler)
		r.Put("/{id}", userHandler.UpdateHandler)
		r.Delete("/{id}", userHandler.DeleteHandler)

		// Rota Especial: Ver apontamentos deste usuário
		r.Get("/{userID}/appointments", appointmentHandler.ListByUser)
	})

	// --- 3. ROTAS DE CONTRATOS (CONTRACTS) ---
	r.Route("/api/contracts", func(r chi.Router) {
		r.Post("/", contractHandler.CreateHandler)
		r.Get("/", contractHandler.GetAllHandler) // Lista Turbinada (com JOIN)
		r.Get("/{id}", contractHandler.GetByIDHandler)
		r.Put("/{id}", contractHandler.UpdateHandler)
		r.Delete("/{id}", contractHandler.DeleteHandler)

		// Rota Especial: Ver apontamentos deste contrato
		r.Get("/{contractID}/appointments", appointmentHandler.ListByContract)
	})

	// --- 4. ROTAS DE APONTAMENTOS (APPOINTMENTS) ---
	r.Route("/api/appointments", func(r chi.Router) {
		r.Post("/", appointmentHandler.CreateHandler) // Lançar horas
		r.Get("/", appointmentHandler.GetAllHandler)  // Visão Admin (Tudo)
		r.Delete("/{id}", appointmentHandler.DeleteHandler)
	})

	return r
}
