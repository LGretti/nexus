package handlers

import "net/http"

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", RootHandler)

	mux.HandleFunc("/empresas/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateEmpresaHandler(w, r)
		case http.MethodGet:
			GetEmpresaByIDHandler(w, r)
		case http.MethodPut:
			UpdateEmpresaHandler(w, r)
		case http.MethodDelete:
			DeleteEmpresaHandler(w, r)
		default:
			RespondWithError(w, http.StatusMethodNotAllowed, "Método inválido")
		}
	})

	mux.HandleFunc("/usuarios/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateUsuarioHandler(w, r)
		case http.MethodGet:
			GetUsuarioByIDHandler(w, r)
		case http.MethodPut:
			UpdateUsuarioHandler(w, r)
		case http.MethodDelete:
			DeleteUsuarioHandler(w, r)
		default:
			RespondWithError(w, http.StatusMethodNotAllowed, "Método inválido")
		}
	})

	return mux
}
