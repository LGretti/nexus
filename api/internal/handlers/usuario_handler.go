package handlers

import (
	"encoding/json"
	"net/http"

	"nexus/api/internal/models"
	"nexus/api/internal/repository"
)

// UsuarioHandler lida com as requisições para usuários.
type UsuarioHandler struct {
	*BaseHandler[*models.Usuario]
	repo repository.UsuarioRepository
}

// NewUsuarioHandler cria um novo handler de usuários, sobrescrevendo o CreateHandler.
func NewUsuarioHandler(repo repository.UsuarioRepository) *UsuarioHandler {
	baseHandler := NewBaseHandler[*models.Usuario](repo, "usuarios")
	handler := &UsuarioHandler{
		BaseHandler: baseHandler,
		repo:        repo,
	}
	handler.CreateHandler = handler.createUsuarioHandler
	return handler
}

// createUsuarioHandler é a implementação customizada para criar um usuário.
func (h *UsuarioHandler) createUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	usuario := h.newModel()
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if usuario.Nome == "" {
		RespondWithError(w, http.StatusBadRequest, "O nome do usuário não pode ser vazio")
		return
	}

	exists, err := h.repo.EmailExists(usuario.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao verificar e-mail")
		return
	}
	if exists {
		RespondWithError(w, http.StatusConflict, "E-mail já cadastrado")
		return
	}

	savedUsuario, err := h.repo.Save(usuario)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao criar usuário")
		return
	}

	RespondWithJSON(w, http.StatusCreated, savedUsuario)
}

// UsuariosRouterHandler delega para o router do handler base.
func (h *UsuarioHandler) UsuariosRouterHandler(w http.ResponseWriter, r *http.Request) {
	h.RouterHandler(w, r)
}
