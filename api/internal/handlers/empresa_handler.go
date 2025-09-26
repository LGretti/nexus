package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"nexus/api/internal/database"
	"nexus/api/internal/models"
)

// Cadastro de Empresa
func CreateEmpresaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	var empresa models.Empresa
	err := json.NewDecoder(r.Body).Decode(&empresa)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	//Validações negociais
	if empresa.Nome == "" {
		RespondWithError(w, http.StatusBadRequest, "O nome da empresa não pode ser vazio")
		return
	}

	if empresa.CNPJ == "" {
		RespondWithError(w, http.StatusBadRequest, "O CNPJ não pode ser vazio")
		return
	}

	//Tudo certo, prepara para salvar
	empresaSalva, err := database.CadEmpresa(empresa)
	if err != nil {
		log.Printf("Erro ao criar empresa: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Erro ao criar empresa")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(empresaSalva)
}

// Cadastro de Empresas em Lote
func CreateEmpresasEmLoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	var novasEmpresas []models.Empresa
	err := json.NewDecoder(r.Body).Decode(&novasEmpresas)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido. Esperado um array de empresas.")
		return
	}

	empresasSalvas, err := database.CadEmpresasEmLote(novasEmpresas)
	if err != nil {
		log.Printf("Erro ao criar empresas em lote: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Erro ao criar empresas em lote")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(empresasSalvas)
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
