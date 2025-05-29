package main

import (
	"arz-synchro-reloj/internal/syncmarca"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var client *syncmarca.Client

func marcacionesHandler(w http.ResponseWriter, r *http.Request) {
	// Parsear la URL para obtener los parámetros posicionales
	// Esperamos: /marcaciones/year/month
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Verificar que tengamos exactamente 3 partes: marcaciones, year, month
	if len(pathParts) != 3 {
		http.Error(w, "URL debe tener el formato: /marcaciones/year/month (ej: /marcaciones/2024/01)", http.StatusBadRequest)
		return
	}

	year := pathParts[1]
	month := pathParts[2]

	// Validar que los parámetros no estén vacíos
	if year == "" || month == "" {
		http.Error(w, "Parámetros year y month son requeridos", http.StatusBadRequest)
		return
	}

	marcaciones, err := client.GetMarcacionesWithAutoConnect(year, month)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error obteniendo marcaciones: %v", err),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(marcaciones)
}

func main() {
	client = syncmarca.NewDefaultClient()

	http.HandleFunc("/marcaciones/", marcacionesHandler)
	fmt.Println("Servidor corriendo en http://localhost:8080")
	fmt.Println("Uso: GET /marcaciones/year/month (ej: /marcaciones/2024/01)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
