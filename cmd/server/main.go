package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"arz-synchro-reloj/internal/syncmarca"
)

func marcacionesHandler(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	month := r.URL.Query().Get("month")
	
	if year == "" || month == "" {
		http.Error(w, "Parámetros year y month son requeridos", http.StatusBadRequest)
		return
	}
	
	client := syncmarca.NewDefaultClient()
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
	http.HandleFunc("/marcaciones", marcacionesHandler)
	http.HandleFunc("/health", healthHandler)
	
	fmt.Println("Servidor corriendo en http://localhost:8080")
	fmt.Println("Endpoints disponibles:")
	fmt.Println("  GET /marcaciones?year=2024&month=04")
	fmt.Println("  GET /health")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
