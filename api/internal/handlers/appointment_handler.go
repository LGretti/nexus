package handlers

import (
	"encoding/json"
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
	baseHandler := NewBaseHandler(repo, "appointments")
	handler := &AppointmentHandler{
		BaseHandler: baseHandler,
		repo:        repo,
	}

	handler.CreateHandler = handler.CreateAppointmentHandler
	handler.GetAllHandler = handler.ListAllWithDetails

	return handler
}

// MÉTODOS BASE CUSTOMIZADOS - Apontar para o Handler

// CreateAppointment godoc
// @Summary      Inicia ou Registra um apontamento
// @Description  Cria um novo registro. Se 'end_time' for omitido ou null, a tarefa inicia "Em Andamento".
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        appointment body models.Appointment true "Dados do Apontamento (EndTime opcional)"
// @Success      201  {object}  models.Appointment
// @Failure      400  {string}  string "Erro de validação"
// @Router       /api/appointments [post]
// 4. Create customizado (caso precise validar horários no futuro)
func (h *AppointmentHandler) CreateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	appt := h.newModel()
	if err := json.NewDecoder(r.Body).Decode(&appt); err != nil {
		RespondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Validação Lógica: Se mandou Data Fim, ela tem que ser maior que Início
	// Se EndTime for nil (nulo), ignora essa validação.
	if appt.EndTime != nil {
		if appt.EndTime.Before(appt.StartTime) {
			RespondWithError(w, http.StatusBadRequest, "A data de fim não pode ser anterior ao início")
			return
		}
	}

	// Grava no banco (O repositório deve estar preparado para aceitar nil)
	savedAppt, err := h.repo.Save(appt)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao salvar apontamento: "+err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, savedAppt)
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

// MÉTODOS ESPECÍFICOS - Apontar para o router

// ListByContract godoc
// @Summary      Lista apontamentos de um contrato
// @Description  Retorna todos os apontamentos vinculados a um ID de contrato específico
// @Tags         contracts
// @Accept       json
// @Produce      json
// @Param        contractID path int true "ID do Contrato"
// @Success      200  {array}  models.Appointment
// @Router       /api/contracts/{contractID}/appointments [get]
// 2. Listar por CONTRATO (Visão Detalhada do Contrato)
func (h *AppointmentHandler) ListAppointmentsByContract(w http.ResponseWriter, r *http.Request) {
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

// ListAppointmentsByUser godoc
// @Summary      Lista apontamentos de um usuário
// @Description  Retorna o histórico de horas de um consultor específico
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userID path int true "ID do Usuário"
// @Success      200  {array}  models.Appointment
// @Failure      400  {string} string "ID inválido"
// @Router       /api/users/{userID}/appointments [get]
// 3. Listar por USUÁRIO (Visão do Consultor)
func (h *AppointmentHandler) ListAppointmentsByUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID do usuário inválido")
		return
	}

	appointments, err := h.repo.GetByUserID(userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao buscar histórico: "+err.Error())
		return
	}

	if appointments == nil {
		RespondWithJSON(w, http.StatusOK, []*models.Appointment{})
		return
	}

	RespondWithJSON(w, http.StatusOK, appointments)
}
