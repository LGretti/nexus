package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"nexus/internal/models"
	"nexus/internal/repository"
	"nexus/internal/utils"
)

// CompanyHandler lida com as requisições para Companies.
type CompanyHandler struct {
	*BaseHandler[*models.Company]
	repo repository.CompanyRepository
}

// NewCompanyHandler cria um novo handler de companies, sobrescrevendo o CreateHandler.
func NewCompanyHandler(repo repository.CompanyRepository) *CompanyHandler {
	baseHandler := NewBaseHandler(repo, "companies")
	handler := &CompanyHandler{
		BaseHandler: baseHandler,
		repo:        repo,
	}
	// Sobrescreve o handler de criação padrão pelo customizado
	handler.CreateHandler = handler.CreateCompanyHandler
	return handler
}

// MÉTODOS BASE CUSTOMIZADOS - Apontar para o Handler

// CreateCompanyHandler lida com a criação de uma ou mais companies.
func (h *CompanyHandler) CreateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var payload interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	var companiesToSave []*models.Company

	if _, ok := payload.([]interface{}); ok {
		jsonBytes, _ := json.Marshal(payload)
		if err := json.Unmarshal(jsonBytes, &companiesToSave); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "JSON de array de companies inválido")
			return
		}
	} else if _, ok := payload.(map[string]interface{}); ok {
		var company *models.Company
		jsonBytes, _ := json.Marshal(payload)
		if err := json.Unmarshal(jsonBytes, &company); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "JSON de company inválido")
			return
		}
		companiesToSave = append(companiesToSave, company)
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "Formato de JSON inválido. Deve ser um objeto ou um array de objetos.")
		return
	}

	for _, company := range companiesToSave {
		if company.Name == "" || company.CNPJ == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Nome e CNPJ são obrigatórios para todas as empresas")
			return
		}
	}

	if len(companiesToSave) == 1 {
		savedCompany, err := h.repo.Save(companiesToSave[0])
		if err != nil {
			if strings.Contains(err.Error(), "companies_cnpj_key") {
				utils.RespondWithError(w, http.StatusConflict, "Este CNPJ já está cadastrado no sistema.")
				return
			}

			utils.RespondWithError(w, http.StatusInternalServerError, "Erro ao criar empresa: "+err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusCreated, savedCompany)
	} else if len(companiesToSave) > 1 {
		savedCompanies, err := h.repo.SaveBatch(companiesToSave)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Erro ao criar empresas em lote: "+err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusCreated, savedCompanies)
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "Nenhuma empresa para cadastrar")
	}
}
