package main

import (
	"arz-synchro-reloj/internal/syncmarca"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
<<<<<<< HEAD
	"strings"
=======
	"arz-synchro-reloj/internal/syncmarca"
>>>>>>> d4855163838ac42b9139da50e497aef86d60bb25
)

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
<<<<<<< HEAD

=======
	
	client := syncmarca.NewDefaultClient()
>>>>>>> d4855163838ac42b9139da50e497aef86d60bb25
	marcaciones, err := client.GetMarcacionesWithAutoConnect(year, month)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error obteniendo marcaciones: %v", err),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(marcaciones)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func main() {
<<<<<<< HEAD
	client = syncmarca.NewDefaultClient()

	http.HandleFunc("/marcaciones/", marcacionesHandler)
	fmt.Println("Servidor corriendo en http://localhost:8080")
	fmt.Println("Uso: GET /marcaciones/year/month (ej: /marcaciones/2024/01)")
=======
	http.HandleFunc("/marcaciones", marcacionesHandler)
	http.HandleFunc("/health", healthHandler)
	
	fmt.Println("Servidor corriendo en http://localhost:8080")
	fmt.Println("Endpoints disponibles:")
	fmt.Println("  GET /marcaciones?year=2024&month=04")
	fmt.Println("  GET /health")
	
>>>>>>> d4855163838ac42b9139da50e497aef86d60bb25
	log.Fatal(http.ListenAndServe(":8080", nil))
}
