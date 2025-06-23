package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	fmt.Println("Iniciando scheduler de sincronización de marcaciones...")
	fmt.Println("El programa se ejecutará todos los días a las 9:00 AM")
	
	// Obtener la ruta del ejecutable del programa principal
	executablePath, err := getExecutablePath()
	if err != nil {
		log.Fatalf("Error obteniendo la ruta del ejecutable: %v", err)
	}
	
	fmt.Printf("Ejecutando: %s\n", executablePath)
	
	// Bucle infinito para ejecutar el programa diariamente
	for {
		// Calcular el próximo tiempo de ejecución (9:00 AM)
		nextRun := getNextRunTime()
		
		// Esperar hasta la próxima ejecución
		waitTime := time.Until(nextRun)
		fmt.Printf("Próxima ejecución programada para: %s (en %v)\n", nextRun.Format("2006-01-02 15:04:05"), waitTime)
		
		time.Sleep(waitTime)
		
		// Ejecutar el programa
		fmt.Printf("Ejecutando sincronización a las %s...\n", time.Now().Format("15:04:05"))
		
		cmd := exec.Command(executablePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		err := cmd.Run()
		if err != nil {
			log.Printf("Error ejecutando el programa: %v", err)
		} else {
			fmt.Println("Sincronización completada exitosamente")
		}
		
		fmt.Println("Esperando hasta la próxima ejecución...")
	}
}

// getExecutablePath obtiene la ruta del ejecutable del programa principal
func getExecutablePath() (string, error) {
	// Obtener el directorio actual
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	
	// Construir la ruta al ejecutable del programa principal
	executablePath := filepath.Join(currentDir, "console")
	
	// Verificar si el archivo existe
	if _, err := os.Stat(executablePath); os.IsNotExist(err) {
		return "", fmt.Errorf("ejecutable no encontrado en: %s", executablePath)
	}
	
	return executablePath, nil
}

// getNextRunTime calcula el próximo tiempo de ejecución a las 9:00 AM
func getNextRunTime() time.Time {
	now := time.Now()
	
	// Crear el tiempo objetivo para hoy a las 9:00 AM
	today := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	
	// Si ya pasó la hora de hoy, programar para mañana
	if now.After(today) {
		today = today.Add(24 * time.Hour)
	}
	
	return today
} 