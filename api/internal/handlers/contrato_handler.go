package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"nexus/api/internal/models"
	"nexus/api/internal/repository"
)

// ContratoHandler lida com as requisições para Contratos.
type ContratoHandler struct {
	*BaseHandler[*models.Contrato]
	repo repository.ContratoRepository
}

// NewContratoHandler cria um novo handler de contratos, sobrescrevendo os handlers.
func NewContratoHandler(repo repository.ContratoRepository) *ContratoHandler {
	baseHandler := NewBaseHandler[*models.Contrato](repo, "contratos")
	handler := &ContratoHandler{
		BaseHandler: baseHandler,
		repo:        repo,
	}
	// Sobrescreve os handlers padrão pelos customizados
	handler.CreateHandler = handler.createContratoHandler
	handler.UpdateHandler = handler.updateContratoHandler
	return handler
}

// ContratosRouterHandler decide qual handler chamar com base na URL.
func (h *ContratoHandler) ContratosRouterHandler(w http.ResponseWriter, r *http.Request) {
	// Rota específica: /empresas/{id}/contratos
	if strings.Contains(r.URL.Path, "/empresas/") {
		h.GetContratosPorEmpresaHandler(w, r)
		return
	}
	// Rotas CRUD genéricas para /contratos/ e /contratos/{id}
	h.RouterHandler(w, r)
}

// createContratoHandler sobrescreve o método base para adicionar validação de data.
func (h *ContratoHandler) createContratoHandler(w http.ResponseWriter, r *http.Request) {
	contrato := h.newModel()
	if err := json.NewDecoder(r.Body).Decode(&contrato); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if contrato.DataFim.Before(contrato.DataInicio) {
		RespondWithError(w, http.StatusBadRequest, "Data de fim não pode ser anterior à data de início")
		return
	}

	savedContrato, err := h.repo.Save(contrato)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao criar contrato: "+err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, savedContrato)
}

// updateContratoHandler sobrescreve o método base para adicionar validação de data.
func (h *ContratoHandler) updateContratoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	contrato := h.newModel()
	if err := json.NewDecoder(r.Body).Decode(&contrato); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if contrato.DataFim.Before(contrato.DataInicio) {
		RespondWithError(w, http.StatusBadRequest, "Data de fim não pode ser anterior à data de início")
		return
	}

	contrato.SetID(id)
	rowsAffected, err := h.repo.Update(contrato)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao atualizar contrato: "+err.Error())
		return
	}
	if rowsAffected == 0 {
		RespondWithError(w, http.StatusNotFound, "Contrato não encontrado")
		return
	}
	RespondWithJSON(w, http.StatusOK, contrato)
}

// GetContratosPorEmpresaHandler lida com a busca de contratos por ID da empresa.
func (h *ContratoHandler) GetContratosPorEmpresaHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 || parts[0] != "empresas" || parts[2] != "contratos" {
		RespondWithError(w, http.StatusBadRequest, "URL mal formada")
		return
	}

	empresaID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID da empresa inválido")
		return
	}

	contratos, err := h.repo.GetPorEmpresaID(empresaID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao buscar contratos por empresa: "+err.Error())
		return
	}
	if len(contratos) == 0 {
		RespondWithJSON(w, http.StatusOK, []*models.Contrato{})
		return
	}
	RespondWithJSON(w, http.StatusOK, contratos)
}
