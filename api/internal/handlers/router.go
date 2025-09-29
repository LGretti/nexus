package handlers

import "net/http"

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", RootHandler)
	mux.HandleFunc("/empresas/", EmpresasRouterHandler)
	mux.HandleFunc("/usuarios/", UsuariosRouterHandler)
	mux.HandleFunc("/contratos/", ContratosRouterHandler)

	return mux
}
