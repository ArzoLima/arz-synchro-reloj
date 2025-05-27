package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"syncmarca"
)

var client *syncmarca.Client

func marcacionesHandler(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	month := r.URL.Query().Get("month")
	
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
	
	http.HandleFunc("/marcaciones", marcacionesHandler)
	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
