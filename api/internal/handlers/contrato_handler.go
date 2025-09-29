package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"nexus/api/internal/database"
	"nexus/api/internal/models"
)

func ContratosRouterHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/contratos/")

	if path == "" { // Rota /contratos/
		switch r.Method {
		case http.MethodPost:
			CreateContratoHandler(w, r)
		// ... outros métodos se aplicável ...
		default:
			RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido para /contratos/")
		}
	} else { // Rota /contratos/{id}
		switch r.Method {
		case http.MethodGet:
			// GetContratoByIDHandler(w, r)
		case http.MethodPut:
			// UpdateContratoHandler(w, r)
		case http.MethodDelete:
			// DeleteContratoHandler(w, r)
		default:
			RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido para /contratos/{id}")
		}
	}
}

// CreateContratoHandler lida com a criação de um novo contrato
func CreateContratoHandler(w http.ResponseWriter, r *http.Request) {
	var novoContrato models.Contrato
	if err := json.NewDecoder(r.Body).Decode(&novoContrato); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	// Validação básica (exemplo)
	if novoContrato.EmpresaID == 0 || novoContrato.HorasContratadas <= 0 {
		RespondWithError(w, http.StatusBadRequest, "empresa_id e horas_contratadas são obrigatórios e devem ser válidos")
		return
	}

	contratoSalvo, err := database.CreateContrato(novoContrato)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao salvar contrato")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contratoSalvo)
}

func GetContratosHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/contratos/")
	var id *int64
	if path != "" {
		parseID, err := strconv.ParseInt(path, 10, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "ID de contrato inválido")
			return
		}
		id = &parseID
	}
	contratos, err := database.GetContratos(id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "erro ao obter lista de contratos: "+err.Error())
		return
	}
	if id != nil && len(contratos) == 0 {
		RespondWithError(w, http.StatusNotFound, "Contrato não encontrado")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contratos)
}
func UpdateContratoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extrair o ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/contratos/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	// 2. Decodificar o corpo da requisição
	var contratoAtualizado models.Contrato
	if err := json.NewDecoder(r.Body).Decode(&contratoAtualizado); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}
	// 3. Atribuir o ID da URL à struct e chamar o banco
	contratoAtualizado.ID = id
	rowsAffected, err := database.UpdateContrato(contratoAtualizado)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao atualizar contrato: "+err.Error())
		return
	}
	// 4. Verificar se alguma linha foi afetada
	if rowsAffected == 0 {
		RespondWithError(w, http.StatusNotFound, "Contrato não encontrado")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func DeleteContratoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extrair o ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/contratos/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// 2. Chamar o banco para deletar o contrato
	rowsAffected, err := database.DeleteContrato(id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao deletar contrato")
		return
	}
	// 3. Verificar se alguma linha foi afetada
	if rowsAffected == 0 {
		RespondWithError(w, http.StatusNotFound, "Contrato não encontrado")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
