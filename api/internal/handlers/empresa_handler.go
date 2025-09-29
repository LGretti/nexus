package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"nexus/api/internal/database"
	"nexus/api/internal/models"
)

func EmpresasRouterHandler(w http.ResponseWriter, r *http.Request) {
	// Divide o path em segmentos: ex: ["empresas", "1", "contratos"]
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Se tem 1 parte (só "empresas")
	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			GetEmpresasHandler(w, r)
		case http.MethodPost:
			CreateEmpresaHandler(w, r)
		default:
			RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido para /empresas/")
		}
		return
	}

	// Se tem 2 partes (ex: "empresas", "1")
	if len(parts) == 2 {
		switch r.Method {
		case http.MethodGet:
			GetEmpresasHandler(w, r)
		case http.MethodPut:
			UpdateEmpresaHandler(w, r)
		case http.MethodDelete:
			DeleteEmpresaHandler(w, r)
		default:
			RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido para /empresas/{id}")
		}
		return
	}

	// Se tem 3 partes e a terceira é "contratos" (ex: "empresas", "1", "contratos")
	if len(parts) == 3 && parts[2] == "contratos" {
		if r.Method == http.MethodGet {
			GetContratosByEmpresaIDHandler(w, r) // Nosso novo handler!
		} else {
			RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido para /empresas/{id}/contratos")
		}
		return
	}

	// Se não caiu em nenhuma regra, rota não encontrada
	RespondWithError(w, http.StatusNotFound, "Rota não encontrada")
}

// Crie o novo handler neste mesmo arquivo
func GetContratosByEmpresaIDHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	empresaID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID de empresa inválido")
		return
	}

	contratos, err := database.GetContratosPorEmpresaID(empresaID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao buscar contratos da empresa")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contratos)
}

// Cadastro de Empresas
func CreateEmpresaHandler(w http.ResponseWriter, r *http.Request) {
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
	path := strings.TrimPrefix(r.URL.Path, "/empresas/")

	var id *int64

	if path != "" {
		parseID, err := strconv.ParseInt(path, 10, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "ID de usuário inválido")
			return
		}
		id = &parseID
	}

	listaEmpresas, err := database.GetEmpresas(id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Erro ao obter lista de usuários: "+err.Error())
		return
	}

	if id != nil && len(listaEmpresas) == 0 {
		RespondWithError(w, http.StatusNotFound, "Empresa não encontrada")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if id != nil {
		json.NewEncoder(w).Encode(listaEmpresas[0])
	} else {
		json.NewEncoder(w).Encode(listaEmpresas)
	}
}

// Atualiza uma empresa
func UpdateEmpresaHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extrair o ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/empresas/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// 2. Decodificar o corpo da requisição
	var empresaAtualizada models.Empresa
	if err := json.NewDecoder(r.Body).Decode(&empresaAtualizada); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	// 3. Atribuir o ID da URL à struct e chamar o banco
	empresaAtualizada.ID = int(id)
	rowsAffected, err := database.UpdateEmpresa(empresaAtualizada)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao atualizar empresa")
		return
	}

	// 4. Verificar se a empresa foi encontrada
	if rowsAffected == 0 {
		RespondWithError(w, http.StatusNotFound, "Empresa não encontrada")
		return
	}

	// 5. Responder com sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(empresaAtualizada)
}

// Deleta uma empresa
func DeleteEmpresaHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extrair o ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/empresas/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// 2. Chamar o banco de dados
	rowsAffected, err := database.DeleteEmpresa(id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao deletar empresa")
		return
	}

	// 3. Verificar se a empresa foi encontrada
	if rowsAffected == 0 {
		RespondWithError(w, http.StatusNotFound, "Empresa não encontrada")
		return
	}

	// 4. Responder com sucesso (sem corpo na resposta)
	w.WriteHeader(http.StatusNoContent)
}
