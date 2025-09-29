package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"nexus/api/internal/database"
	"nexus/api/internal/models"
)

func UsuariosRouterHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/usuarios/")

	// Se o path está vazio, a rota é /usuarios/
	if path == "" {
		switch r.Method {
		case http.MethodGet:
			GetUsuariosHandler(w, r) // Lista todos os usuários
		case http.MethodPost:
			CreateUsuarioHandler(w, r) // Cria usuário(s)
		default:
			RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido para /usuarios/")
		}
	} else { // Se o path não está vazio, ele deve ser um ID
		switch r.Method {
		case http.MethodGet:
			GetUsuariosHandler(w, r)
		case http.MethodPut:
			UpdateUsuarioHandler(w, r)
		case http.MethodDelete:
			DeleteUsuarioHandler(w, r)
		default:
			RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido para /usuarios/{id}")
		}
	}
}

func CreateUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	var usuario models.Usuario
	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if usuario.Nome == "" {
		RespondWithError(w, http.StatusBadRequest, "O nome do usuário não pode ser vazio")
		return
	}

	usuarioSalvo, err := database.CreateUsuario(usuario)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao criar usuário")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuarioSalvo)
}

func GetUsuariosHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/usuarios/")

	var id *int64

	if path != "" {
		parseID, err := strconv.ParseInt(path, 10, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "ID de usuário inválido")
			return
		}
		id = &parseID
	}

	listaUsuarios, err := database.GetUsuarios(id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "erro ao obter lista de usuários: "+err.Error())
		return
	}

	if id != nil && len(listaUsuarios) == 0 {
		RespondWithError(w, http.StatusNotFound, "Empresa não encontrada")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if id != nil {
		json.NewEncoder(w).Encode(listaUsuarios[0])
	} else {
		json.NewEncoder(w).Encode(listaUsuarios)
	}
}

func UpdateUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extrair o ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/usuarios/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// 2. Decodificar o corpo da requisição
	var usuarioAtualizado models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuarioAtualizado); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	// 3. Atribuir o ID da URL à struct e chamar o banco
	usuarioAtualizado.ID = int(id)
	rowsAffected, err := database.UpdateUsuario(usuarioAtualizado)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao atualizar usuário")
		return
	}

	// 4. Verificar se o usuário foi encontrado
	if rowsAffected == 0 {
		RespondWithError(w, http.StatusNotFound, "Usuário não encontrado")
		return
	}

	// 5. Responder com sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usuarioAtualizado)
}

func DeleteUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extrair o ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/usuarios/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// 2. Chamar o banco de dados
	rowsAffected, err := database.DeleteUsuario(id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao deletar usuário")
		return
	}

	// 3. Verificar se o usuário foi encontrado
	if rowsAffected == 0 {
		RespondWithError(w, http.StatusNotFound, "Usuário não encontrado")
		return
	}

	// 4. Responder com sucesso (sem corpo na resposta)
	w.WriteHeader(http.StatusNoContent)
}
