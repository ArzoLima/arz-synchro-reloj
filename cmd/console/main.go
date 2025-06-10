package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"arz-synchro-reloj/internal/syncmarca"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type User struct {
	ID       int
	Username string
}

// Marcacion representa cada registro extraído de la base de datos.
type Marcas struct {
	id_trabajador string
	Fecha         string
	Hora          string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Conexión exitosa a MySQL desde Go")

	rows, err := db.Query("SELECT id, username FROM arz_intranet.auth_user;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Solicitar parámetros
	fmt.Print("Ingresa el año (YYYY): ")
	var year string
	fmt.Scanln(&year)

	fmt.Print("Ingresa el mes (MM): ")
	var month string
	fmt.Scanln(&month)

	// Crear cliente y obtener marcaciones
	client := syncmarca.NewDefaultClient()
	marcaciones, err := client.GetMarcacionesWithAutoConnect(year, month)
	if err != nil {
		log.Fatalf("Error obteniendo marcaciones: %v", err)
	}

	var marcasProcesadas []Marcas
	for _, marcacion := range marcaciones {
		for _, user := range users {
			// Extraer los últimos 8 dígitos del IdEmpleado
			if len(marcacion.IdEmpleado) >= 8 {
				idSuffix := marcacion.IdEmpleado[len(marcacion.IdEmpleado)-8:]
				if idSuffix == user.Username {
					marcasProcesadas = append(marcasProcesadas, Marcas{
						id_trabajador: user.Username,
						Fecha:         marcacion.Fecha,
						Hora:          marcacion.Hora,
					})
				}
			}
		}
	}

	// Imprimir las marcaciones procesadas
	fmt.Println("Marcaciones procesadas:")
	for _, m := range marcasProcesadas {
		fmt.Printf("id_trabajador: %s, Fecha: %s, Hora: %s\n", m.id_trabajador, m.Fecha, m.Hora)
	}
}