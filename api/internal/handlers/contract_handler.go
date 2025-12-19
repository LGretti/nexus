package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"nexus/internal/models"
	"nexus/internal/repository"
	"nexus/internal/utils"

	"github.com/go-chi/chi/v5"
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
	handler.CreateHandler = handler.CreateContractHandler
	handler.UpdateHandler = handler.UpdateContractHandler
	handler.GetAllHandler = handler.ListContracts
	return handler
}

// MÉTODOS BASE CUSTOMIZADOS - Apontar para o Handler

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
// CreateContractHandler sobrescreve o método base para adicionar validação de data.
func (h *ContractHandler) CreateContractHandler(w http.ResponseWriter, r *http.Request) {
	contract := h.newModel()
	if err := json.NewDecoder(r.Body).Decode(&contract); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if contract.EndDate.Before(contract.StartDate) {
		utils.RespondWithError(w, http.StatusBadRequest, "Data de fim não pode ser anterior à data de início")
		return
	}

	savedContract, err := h.repo.Save(contract)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Erro ao criar contrato: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, savedContract)
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
// UpdateContractHandler sobrescreve o método base para adicionar validação de data.
func (h *ContractHandler) UpdateContractHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	contract := h.newModel()
	if err := json.NewDecoder(r.Body).Decode(&contract); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if contract.EndDate.Before(contract.StartDate) {
		utils.RespondWithError(w, http.StatusBadRequest, "Data de fim não pode ser anterior à data de início")
		return
	}

	contract.SetID(id)
	rowsAffected, err := h.repo.Update(contract)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Erro ao atualizar contrato: "+err.Error())
		return
	}
	if rowsAffected == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Contrato não encontrado")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, contract)
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
		utils.RespondWithError(w, http.StatusInternalServerError, "Erro ao buscar contratos: "+err.Error())
		return
	}

	// Se não tiver nada, retorna array vazio [] em vez de null
	if contracts == nil {
		utils.RespondWithJSON(w, http.StatusOK, []*models.Contract{})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, contracts)
}

// MÉTODOS ESPECÍFICOS - Apontar para o router

// ListContractsByCompany lida com a busca de contratos por ID da empresa.
func (h *ContractHandler) ListContractsByCompany(w http.ResponseWriter, r *http.Request) {
	// 1. O Chi já separou o ID pra gente. É só pegar.
	companyIDStr := chi.URLParam(r, "companyID")

	// Se por acaso vier vazio (caso você teste sem o router), validamos:
	if companyIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "ID da empresa não fornecido na URL")
		return
	}

	// 2. Conversão segura
	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "ID da empresa inválido")
		return
	}

	// 3. Chamada ao Banco (sem mexer na lógica)
	contracts, err := h.repo.GetByCompanyID(companyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Erro ao buscar contratos: "+err.Error())
		return
	}

	// 4. Retorno Vazio (Array vazio é melhor que null)
	if contracts == nil {
		utils.RespondWithJSON(w, http.StatusOK, []*models.Contract{})
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, contracts)
}
