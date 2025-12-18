package handlers

import (
	"encoding/json"
	"net/http"

	"nexus/api/internal/models"
	"nexus/api/internal/repository"
)

// UserHandler lida com as requisições para usuários.
type UserHandler struct {
	*BaseHandler[*models.User]
	repo repository.UserRepository
}

// NewUserHandler cria um novo handler de usuários, sobrescrevendo o CreateHandler.
func NewUserHandler(repo repository.UserRepository) *UserHandler {
	baseHandler := NewBaseHandler(repo, "users")
	handler := &UserHandler{
		BaseHandler: baseHandler,
		repo:        repo,
	}
	handler.CreateHandler = handler.createUserHandler
	return handler
}

// createUserHandler é a implementação customizada para criar um usuário.
func (h *UserHandler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	user := h.newModel()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if user.Name == "" {
		RespondWithError(w, http.StatusBadRequest, "O nome do usuário não pode ser vazio")
		return
	}

	exists, err := h.repo.EmailExists(user.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao verificar e-mail")
		return
	}
	if exists {
		RespondWithError(w, http.StatusConflict, "E-mail já cadastrado")
		return
	}

	savedUser, err := h.repo.Save(user)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao criar usuário")
		return
	}

	RespondWithJSON(w, http.StatusCreated, savedUser)
}

// UserHandler delega para o router do handler base.
func (h *UserHandler) UsersRouterHandler(w http.ResponseWriter, r *http.Request) {
	h.RouterHandler(w, r)
}
