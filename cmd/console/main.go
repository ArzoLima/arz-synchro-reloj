package main

import (
	"fmt"
	"log"
	"syncmarca"
)

func main() {
	// Opción 1: Usar cliente por defecto
	client := syncmarca.NewDefaultClient()
	
	// Opción 2: Usar configuración personalizada
	// client := syncmarca.NewClient(syncmarca.Config{
	//     Server:                "MISERVIDOR",
	//     Database:              "MiBaseDatos",
	//     UseIntegratedSecurity: true,
	//     Encrypt:               false,
	// })

	// Obtener marcaciones con conexión automática
	marcaciones, err := client.GetMarcacionesWithAutoConnect("2024", "04")
	if err != nil {
		log.Fatalf("Error obteniendo marcaciones: %v", err)
	}

	// Mostrar resultados
	fmt.Printf("Encontradas %d marcaciones:\n", len(marcaciones))
	for _, m := range marcaciones {
		fmt.Printf("Empleado: %s, Fecha: %s, Hora: %s\n", 
			m.IdEmpleado, m.Fecha, m.Hora)
	}
}
