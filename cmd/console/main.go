package main
//recupercion de version mysql 8.0.33
import (
<<<<<<< HEAD
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
=======
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
>>>>>>> d4855163838ac42b9139da50e497aef86d60bb25
	marcaciones, err := client.GetMarcacionesWithAutoConnect(year, month)
	if err != nil {
		log.Fatalf("Error obteniendo marcaciones: %v", err)
	}

<<<<<<< HEAD
	// Mostrar resultados
	fmt.Printf("Encontradas %d marcaciones para %s/%s:\n", len(marcaciones), month, year)
	for _, m := range marcaciones {
		fmt.Printf("Empleado: %s, Fecha: %s, Hora: %s\n",
			m.IdEmpleado, m.Fecha, m.Hora)
=======
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
>>>>>>> d4855163838ac42b9139da50e497aef86d60bb25
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