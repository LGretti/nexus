package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"nexus/api/internal/database"
	"nexus/api/internal/models"
)

func CreateUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

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

	usuarioSalvo, err := database.CadUsuario(usuario)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao criar usuário")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuarioSalvo)
}

func GetUsuariosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	usuarios, err := database.GetUsuarios()
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "erro ao obter lista de usuários: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usuarios)
}

func GetUsuarioByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	idStr := r.URL.Path[len("/usuarios/"):]
	if idStr == "" {
		GetUsuariosHandler(w, r)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID do usuário inválido")
		return
	}

	usuario, err := database.GetUsuarioByID(id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Erro ao obter usuário: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usuario)
}

func UpdateUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	idStr := r.URL.Path[len("/usuarios/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID do usuário inválido")
		return
	}

	usuario, err := database.GetUsuarioByID(id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Erro ao localizar usuário: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usuario)
}

func DeleteUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	idStr := r.URL.Path[len("/usuarios/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID do usuário não localizado")
		return
	}

	err = database.DeleteUsuario(id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
