package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"nexus/api/internal/models"
	"nexus/api/internal/repository"
)

// ContractHandler lida com as requisições para Contratos.
type ContractHandler struct {
	*BaseHandler[*models.Contract]
	repo repository.ContractRepository
}

// NewContractHandler cria um novo handler de contratos, sobrescrevendo os handlers.
func NewContractHandler(repo repository.ContractRepository) *ContractHandler {
	baseHandler := NewBaseHandler(repo, "contracts")
	handler := &ContractHandler{
		BaseHandler: baseHandler,
		repo:        repo,
	}
	// Sobrescreve os handlers padrão pelos customizados
	handler.CreateHandler = handler.createContractHandler
	handler.UpdateHandler = handler.updateContractHandler
	handler.GetAllHandler = handler.ListContracts
	return handler
}

// ContractsRouterHandler decide qual handler chamar com base na URL.
func (h *ContractHandler) ContractsRouterHandler(w http.ResponseWriter, r *http.Request) {
	// Rota específica: /companies/{id}/contracts
	if strings.Contains(r.URL.Path, "/companies/") {
		h.GetContractsByCompanyHandler(w, r)
		return
	}
	// Rotas CRUD genéricas para /contracts/ e /contracts/{id}
	h.RouterHandler(w, r)
}

// createContractHandler godoc
// @Summary      Cria um novo contrato
// @Description  Valida as datas e insere um novo contrato no banco
// @Tags         contracts
// @Accept       json
// @Produce      json
// @Param        contract body models.Contract true "Objeto Contrato"
// @Success      201  {object}  models.Contract
// @Failure      400  {string}  string "Erro de validação"
// @Router       /api/contracts [post]
// createContractHandler sobrescreve o método base para adicionar validação de data.
func (h *ContractHandler) createContractHandler(w http.ResponseWriter, r *http.Request) {
	contract := h.newModel()
	if err := json.NewDecoder(r.Body).Decode(&contract); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if contract.EndDate.Before(contract.StartDate) {
		RespondWithError(w, http.StatusBadRequest, "Data de fim não pode ser anterior à data de início")
		return
	}

	savedContract, err := h.repo.Save(contract)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao criar contrato: "+err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, savedContract)
}

// updateContractHandler godoc
// @Summary      Atualiza um contrato existente
// @Description  Atualiza os dados de um contrato pelo ID
// @Tags         contracts
// @Accept       json
// @Produce      json
// @Param        id   path      int             true "ID do Contrato"
// @Param        contract body models.Contract true "Objeto Contrato Atualizado"
// @Success      200  {object}  models.Contract
// @Failure      400  {string}  string "ID inválido ou datas incorretas"
// @Failure      404  {string}  string "Contrato não encontrado"
// @Router       /api/contracts/{id} [put]
// updateContractHandler sobrescreve o método base para adicionar validação de data.
func (h *ContractHandler) updateContractHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	contract := h.newModel()
	if err := json.NewDecoder(r.Body).Decode(&contract); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if contract.EndDate.Before(contract.StartDate) {
		RespondWithError(w, http.StatusBadRequest, "Data de fim não pode ser anterior à data de início")
		return
	}

	contract.SetID(id)
	rowsAffected, err := h.repo.Update(contract)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao atualizar contrato: "+err.Error())
		return
	}
	if rowsAffected == 0 {
		RespondWithError(w, http.StatusNotFound, "Contrato não encontrado")
		return
	}
	RespondWithJSON(w, http.StatusOK, contract)
}

// GetContractsByCompanyHandler lida com a busca de contratos por ID da empresa.
func (h *ContractHandler) GetContractsByCompanyHandler(w http.ResponseWriter, r *http.Request) {
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

	contracts, err := h.repo.GetByCompanyID(empresaID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao buscar contratos por empresa: "+err.Error())
		return
	}
	if len(contracts) == 0 {
		RespondWithJSON(w, http.StatusOK, []*models.Contract{})
		return
	}
	RespondWithJSON(w, http.StatusOK, contracts)
}

// ListContracts godoc
// @Summary      Lista todos os contratos
// @Description  Retorna a lista completa de contratos com dados da empresa (JOIN)
// @Tags         contracts
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Contract
// @Router       /api/contracts [get]
func (h *ContractHandler) ListContracts(w http.ResponseWriter, r *http.Request) {
	contracts, err := h.repo.GetAllWithCompany()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao buscar contratos: "+err.Error())
		return
	}

	// Se não tiver nada, retorna array vazio [] em vez de null
	if contracts == nil {
		RespondWithJSON(w, http.StatusOK, []*models.Contract{})
		return
	}

	RespondWithJSON(w, http.StatusOK, contracts)
}
