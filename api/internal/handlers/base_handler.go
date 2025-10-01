package handlers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"nexus/api/internal/models"
	"nexus/api/internal/repository"
)

// BaseHandler é um handler genérico para operações CRUD.
// Usa campos de função para permitir a sobrescrita de comportamento.
type BaseHandler[T models.Model] struct {
	repo           repository.Repository[T]
	routeName      string
	CreateHandler  http.HandlerFunc
	GetAllHandler  http.HandlerFunc
	GetByIDHandler http.HandlerFunc
	UpdateHandler  http.HandlerFunc
	DeleteHandler  http.HandlerFunc
}

// NewBaseHandler cria uma nova instância de BaseHandler com handlers padrão.
func NewBaseHandler[T models.Model](repo repository.Repository[T], routeName string) *BaseHandler[T] {
	h := &BaseHandler[T]{
		repo:      repo,
		routeName: routeName,
	}
	// Define os handlers padrão
	h.CreateHandler = h.createHandlerDefault
	h.GetAllHandler = h.getAllHandlerDefault
	h.GetByIDHandler = h.getByIDHandlerDefault
	h.UpdateHandler = h.updateHandlerDefault
	h.DeleteHandler = h.deleteHandlerDefault
	return h
}

// RouterHandler delega a requisição para o handler apropriado com base no método HTTP e no path.
func (h *BaseHandler[T]) RouterHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/"+h.routeName+"/")

	if path == "" || path == "/" {
		switch r.Method {
		case http.MethodGet:
			h.GetAllHandler(w, r)
		case http.MethodPost:
			h.CreateHandler(w, r)
		default:
			RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido para /"+h.routeName)
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			h.GetByIDHandler(w, r)
		case http.MethodPut:
			h.UpdateHandler(w, r)
		case http.MethodDelete:
			h.DeleteHandler(w, r)
		default:
			RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido para /"+h.routeName+"/{id}")
		}
	}
}

func (h *BaseHandler[T]) newModel() T {
	var t T
	return reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
}

// Implementações padrão dos handlers
func (h *BaseHandler[T]) createHandlerDefault(w http.ResponseWriter, r *http.Request) {
	model := h.newModel()
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}
	savedModel, err := h.repo.Save(model)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao criar "+h.routeName+": "+err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, savedModel)
}

func (h *BaseHandler[T]) getAllHandlerDefault(w http.ResponseWriter, r *http.Request) {
	models, err := h.repo.Get(nil)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao obter "+h.routeName+": "+err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, models)
}

func (h *BaseHandler[T]) getByIDHandlerDefault(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	models, err := h.repo.Get(&id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao buscar "+h.routeName+": "+err.Error())
		return
	}
	if len(models) == 0 {
		RespondWithError(w, http.StatusNotFound, h.routeName+" não encontrado")
		return
	}
	RespondWithJSON(w, http.StatusOK, models[0])
}

func (h *BaseHandler[T]) updateHandlerDefault(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	model := h.newModel()
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}
	model.SetID(id)
	rowsAffected, err := h.repo.Update(model)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao atualizar "+h.routeName+": "+err.Error())
		return
	}
	if rowsAffected == 0 {
		RespondWithError(w, http.StatusNotFound, h.routeName+" não encontrado")
		return
	}
	RespondWithJSON(w, http.StatusOK, model)
}

func (h *BaseHandler[T]) deleteHandlerDefault(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	rowsAffected, err := h.repo.Delete(id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao deletar "+h.routeName+": "+err.Error())
		return
	}
	if rowsAffected == 0 {
		RespondWithError(w, http.StatusNotFound, h.routeName+" não encontrado")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BaseHandler[T]) parseID(r *http.Request) (int64, error) {
	path := strings.TrimPrefix(r.URL.Path, "/"+h.routeName+"/")
	return strconv.ParseInt(path, 10, 64)
}
