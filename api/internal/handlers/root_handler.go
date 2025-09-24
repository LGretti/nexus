package handlers

import (
	"fmt"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	// Escreve uma resposta simples de volta para o cliente.
	fmt.Fprintf(w, "API do Nexus está no ar! (v1)")
}