package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"arz-synchro-reloj/internal/syncmarca"
)

func main() {
	// Definir flags para año y mes
	var year string
	var month string

	// Obtener año y mes actual como valores por defecto
	currentTime := time.Now()
	defaultYear := strconv.Itoa(currentTime.Year())
	defaultMonth := fmt.Sprintf("%02d", int(currentTime.Month()))

	flag.StringVar(&year, "year", defaultYear, "Año de las marcaciones (formato: 2024)")
	flag.StringVar(&year, "y", defaultYear, "Año de las marcaciones (formato: 2024) - versión corta")
	flag.StringVar(&month, "month", defaultMonth, "Mes de las marcaciones (formato: 01-12)")
	flag.StringVar(&month, "m", defaultMonth, "Mes de las marcaciones (formato: 01-12) - versión corta")

	flag.Parse()

	// Validar formato del mes (debe ser 01-12)
	if len(month) == 1 {
		month = "0" + month
	}

	// Validar que el mes esté en rango válido
	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		log.Fatalf("Error: El mes debe ser un número entre 1 y 12. Valor recibido: %s", month)
	}

	// Validar formato del año
	yearInt, err := strconv.Atoi(year)
	if err != nil || yearInt < 1900 || yearInt > 2100 {
		log.Fatalf("Error: El año debe ser un número válido entre 1900 y 2100. Valor recibido: %s", year)
	}

	fmt.Printf("Obteniendo marcaciones para: %s/%s\n", month, year)

	// Opción 1: Usar cliente por defecto
	client := syncmarca.NewDefaultClient()

	// Opción 2: Usar configuración personalizada
	// client := syncmarca.NewClient(syncmarca.Config{
	//     Server:                "MISERVIDOR",
	//     Database:              "MiBaseDatos",
	//     UseIntegratedSecurity: true,
	//     Encrypt:               false,
	// })

	// Obtener marcaciones con conexión automática usando los parámetros
	marcaciones, err := client.GetMarcacionesWithAutoConnect(year, month)
	if err != nil {
		log.Fatalf("Error obteniendo marcaciones: %v", err)
	}

	// Mostrar resultados
	fmt.Printf("Encontradas %d marcaciones para %s/%s:\n", len(marcaciones), month, year)
	for _, m := range marcaciones {
		fmt.Printf("Empleado: %s, Fecha: %s, Hora: %s\n",
			m.IdEmpleado, m.Fecha, m.Hora)
	}
}
