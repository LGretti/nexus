package handlers

import (
	"net/http"
)

func NewRouter(empresaHandler *EmpresaHandler, usuarioHandler *UsuarioHandler, contratoHandler *ContratoHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", RootHandler)
	mux.HandleFunc("/empresas/", empresaHandler.EmpresasRouterHandler)
	mux.HandleFunc("/usuarios/", usuarioHandler.UsuariosRouterHandler)
	mux.HandleFunc("/contratos/", contratoHandler.ContratosRouterHandler)

	return mux
}
