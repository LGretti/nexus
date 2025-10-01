package handlers

import (
	"encoding/json"
	"net/http"

	"nexus/api/internal/models"
	"nexus/api/internal/repository"
)

// EmpresaHandler lida com as requisições para Empresas.
type EmpresaHandler struct {
	*BaseHandler[*models.Empresa]
	repo repository.EmpresaRepository
}

// NewEmpresaHandler cria um novo handler de empresas, sobrescrevendo o CreateHandler.
func NewEmpresaHandler(repo repository.EmpresaRepository) *EmpresaHandler {
	baseHandler := NewBaseHandler[*models.Empresa](repo, "empresas")
	handler := &EmpresaHandler{
		BaseHandler: baseHandler,
		repo:        repo,
	}
	// Sobrescreve o handler de criação padrão pelo customizado
	handler.CreateHandler = handler.createEmpresaHandler
	return handler
}

// createEmpresaHandler lida com a criação de uma ou mais empresas.
func (h *EmpresaHandler) createEmpresaHandler(w http.ResponseWriter, r *http.Request) {
	var payload interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	var empresasParaSalvar []*models.Empresa

	if _, ok := payload.([]interface{}); ok {
		jsonBytes, _ := json.Marshal(payload)
		if err := json.Unmarshal(jsonBytes, &empresasParaSalvar); err != nil {
			RespondWithError(w, http.StatusBadRequest, "JSON de array de empresas inválido")
			return
		}
	} else if _, ok := payload.(map[string]interface{}); ok {
		var empresa *models.Empresa
		jsonBytes, _ := json.Marshal(payload)
		if err := json.Unmarshal(jsonBytes, &empresa); err != nil {
			RespondWithError(w, http.StatusBadRequest, "JSON de empresa inválido")
			return
		}
		empresasParaSalvar = append(empresasParaSalvar, empresa)
	} else {
		RespondWithError(w, http.StatusBadRequest, "Formato de JSON inválido. Deve ser um objeto ou um array de objetos.")
		return
	}

	for _, emp := range empresasParaSalvar {
		if emp.Nome == "" || emp.CNPJ == "" {
			RespondWithError(w, http.StatusBadRequest, "Nome e CNPJ são obrigatórios para todas as empresas")
			return
		}
	}

	if len(empresasParaSalvar) == 1 {
		savedEmpresa, err := h.repo.Save(empresasParaSalvar[0])
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Erro ao criar empresa: "+err.Error())
			return
		}
		RespondWithJSON(w, http.StatusCreated, savedEmpresa)
	} else if len(empresasParaSalvar) > 1 {
		savedEmpresas, err := h.repo.SaveBatch(empresasParaSalvar)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Erro ao criar empresas em lote: "+err.Error())
			return
		}
		RespondWithJSON(w, http.StatusCreated, savedEmpresas)
	} else {
		RespondWithError(w, http.StatusBadRequest, "Nenhuma empresa para cadastrar")
	}
}

// EmpresasRouterHandler delega para o router do handler base.
func (h *EmpresaHandler) EmpresasRouterHandler(w http.ResponseWriter, r *http.Request) {
	h.RouterHandler(w, r)
}
