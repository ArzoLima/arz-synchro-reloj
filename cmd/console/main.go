package main
//recupercion de version mysql 8.0.33
import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"arz-synchro-reloj/internal/syncmarca"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type User struct {
	ID       int
	Username string
}

// Marcas representa los datos procesados que se insertarán en la base de datos.
type Marcas struct {
	trabajador_id int
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

	// Obtener año y mes actual
	currentTime := time.Now()
	year := strconv.Itoa(currentTime.Year())
	month := fmt.Sprintf("%02d", int(currentTime.Month()))
	fmt.Printf("Procesando marcaciones para: %s/%s\n", month, year)

	// Crear cliente y obtener marcaciones
	client := syncmarca.NewDefaultClient()
	marcaciones, err := client.GetMarcacionesWithAutoConnect(year, month)
	if err != nil {
		log.Fatalf("Error obteniendo marcaciones: %v", err)
	}

	var marcasProcesadas []Marcas
	for _, marcacion := range marcaciones {
		for _, user := range users {
			if len(marcacion.IdEmpleado) >= 8 {
				idSuffix := marcacion.IdEmpleado[len(marcacion.IdEmpleado)-8:]
				if idSuffix == user.Username {
					marcasProcesadas = append(marcasProcesadas, Marcas{
						trabajador_id: user.ID,
						Fecha:         marcacion.Fecha,
						Hora:          marcacion.Hora,
					})
				}
			}
		}
	}

	// Iniciar transacción
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error al iniciar la transacción: %v", err)
	}

	// Borrar registros existentes para el mes y año, usando las funciones YEAR() y MONTH() para mayor precisión.
	_, err = tx.Exec("DELETE FROM arz_intranet.permisos_marcassynchro WHERE YEAR(fecha) = ? AND MONTH(fecha) = ?", year, month)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Error al borrar registros existentes: %v", err)
	}

	// Insertar nuevas marcaciones
	stmt, err := tx.Prepare("INSERT INTO arz_intranet.permisos_marcassynchro(trabajador_id, fecha, hora, fecha_synchro) VALUES(?, ?, ?, NOW())")
	if err != nil {
		tx.Rollback()
		log.Fatalf("Error al preparar la inserción: %v", err)
	}
	defer stmt.Close()

	for _, m := range marcasProcesadas {
		_, err := stmt.Exec(m.trabajador_id, m.Fecha, m.Hora)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Error al insertar marcación: %v", err)
		}
	}

	// Confirmar transacción
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error al confirmar la transacción: %v", err)
	}

	fmt.Println("Sincronización completada con éxito.")
}