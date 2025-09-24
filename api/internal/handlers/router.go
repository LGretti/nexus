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
			RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		}
	})

	mux.HandleFunc("/empresas/lote", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateEmpresasEmLoteHandler(w, r)
		default:
			RespondWithError(w, http.StatusMethodNotAllowed, "Método não permitido")
		}
	})

	return mux
}
