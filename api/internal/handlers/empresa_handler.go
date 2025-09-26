package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"nexus/api/internal/database"
	"nexus/api/internal/models"
)

// Cadastro de Empresas
func CreateEmpresaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	var payload interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	var empresasParaSalvar []models.Empresa

	switch p := payload.(type) {
	case map[string]interface{}:
		var empresa models.Empresa
		jsonBytes, _ := json.Marshal(p)
		if err := json.Unmarshal(jsonBytes, &empresa); err != nil {
			RespondWithError(w, http.StatusBadRequest, "JSON de empresa inválido")
			return
		}
		empresasParaSalvar = append(empresasParaSalvar, empresa)

	case []interface{}:
		jsonBytes, _ := json.Marshal(p)
		if err := json.Unmarshal(jsonBytes, &empresasParaSalvar); err != nil {
			RespondWithError(w, http.StatusBadRequest, "JSON de empresas inválido")
			return
		}

	default:
		RespondWithError(w, http.StatusBadRequest, "Formato de JSON inválido. Deve ser um objeto ou um array de objetos.")
		return
	}

	if len(empresasParaSalvar) == 0 {
		RespondWithError(w, http.StatusBadRequest, "Nenhuma empresa para cadastrar")
		return
	}
	for _, emp := range empresasParaSalvar {
		if emp.Nome == "" || emp.CNPJ == "" {
			RespondWithError(w, http.StatusBadRequest, "O Nome e CNPJ são obrigatórios!")
			return
		}
	}

	empresasSalvas, err := database.CreateEmpresasBatch(empresasParaSalvar)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao criar empresa(s): "+err.Error())
		return
	}

	log.Printf("%d empresa(s) criada(s) com sucesso.\n", len(empresasSalvas))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if len(empresasSalvas) > 1 {
		json.NewEncoder(w).Encode(empresasSalvas)
		return
	}
	json.NewEncoder(w).Encode(empresasSalvas[0])
}

// Retorna todas as empresas cadastradas
func GetEmpresasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	empresas, err := database.GetEmpresas()
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Erro ao obter lista de empresas: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(empresas)
}

// Retorna uma empresa pelo ID
func GetEmpresaByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	// Extrai o ID da URL. Ex: /empresas/1
	idStr := r.URL.Path[len("/empresas/"):]
	if idStr == "" {
		// Se não houver ID, chama o handler que lista todas as empresas.
		GetEmpresasHandler(w, r)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID da empresa inválido")
		return
	}

	empresa, err := database.GetEmpresaByID(id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Erro ao obter empresa: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(empresa)
}

// Atualiza uma empresa
func UpdateEmpresaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	var empresaAtualizada models.Empresa
	err := json.NewDecoder(r.Body).Decode(&empresaAtualizada)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	idStr := r.URL.Path[len("/empresas/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID da empresa inválido")
		return
	}

	empresa, err := database.UpdateEmpresa(id, empresaAtualizada)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Erro ao localizar empresa: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(empresa)
}

// Deleta uma empresa
func DeleteEmpresaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	idStr := r.URL.Path[len("/empresas/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID da empresa não localizado")
		return
	}

	err = database.DeleteEmpresa(id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
