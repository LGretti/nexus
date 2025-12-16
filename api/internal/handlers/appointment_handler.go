package handlers

import (
	"net/http"
	"strconv"

	"nexus/api/internal/models"
	"nexus/api/internal/repository"

	"github.com/go-chi/chi/v5"
)

type AppointmentHandler struct {
	*BaseHandler[*models.Appointment]
	repo repository.AppointmentRepository
}

func NewAppointmentHandler(repo repository.AppointmentRepository) *AppointmentHandler {
	baseHandler := NewBaseHandler[*models.Appointment](repo, "appointments")
	handler := &AppointmentHandler{
		BaseHandler: baseHandler,
		repo:        repo,
	}

	handler.CreateHandler = handler.CreateAppointmentHandler
	handler.GetAllHandler = handler.ListAllWithDetails

	return handler
}

// AppontmentsRouterHandler decide qual handler chamar com base na URL.
func (h *AppointmentHandler) ListAllWithDetails(w http.ResponseWriter, r *http.Request) {
	appointments, err := h.repo.GetAllWithContract()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, appointments)
}

// 2. Listar por CONTRATO (Visão Detalhada do Contrato)
func (h *AppointmentHandler) ListByContract(w http.ResponseWriter, r *http.Request) {
	// Pega o {contractID} da URL usando o Chi
	idStr := chi.URLParam(r, "contractID")
	contractID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID do contrato inválido")
		return
	}

	appointments, err := h.repo.GetByContractID(contractID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, appointments)
}

// 3. Listar por USUÁRIO (Visão do Consultor)
func (h *AppointmentHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID do usuário inválido")
		return
	}

	appointments, err := h.repo.GetByUserID(userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, appointments)
}

// 4. Create customizado (caso precise validar horários no futuro)
func (h *AppointmentHandler) CreateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	// Por enquanto, usa o padrão do BaseHandler, mas chamando sua lógica se precisar
	h.BaseHandler.createHandlerDefault(w, r)
}
