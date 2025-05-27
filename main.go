package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

// Marcacion representa cada registro extraído de la base de datos.
type Marcacion struct {
	IdEmpleado string `json:"IdEmpleado"`
	Fecha      string `json:"Fecha"`
	Hora       string `json:"Hora"`
}

var db *sql.DB // Variable global para la conexión a la base de datos

func main() {
	// Configura la cadena de conexión.
	// Usando autenticación integrada (Windows Authentication) y la instancia predeterminada en "SVRDP".
	dsn := "server=SVRDP;" +
		"database=ZKTime;" +
		"integrated security=true;" +
		"encrypt=disable"

	var err error
	db, err = sql.Open("sqlserver", dsn)
	if err != nil {
		log.Fatal("Error al crear la conexión: ", err)
	}
	defer db.Close()

	// Prueba la conexión con Ping.
	err = db.Ping()
	if err != nil {
		log.Fatal("Error conectando a la base de datos: ", err)
	}
	fmt.Println("Conexión a SQL Server establecida exitosamente")

	// Configura el handler para solicitudes HTTP.
	// Se espera que la URL sea de la forma /{año}/{mes}
	http.HandleFunc("/", queryHandler)

	// Inicia el servidor en el puerto 8080.
	log.Println("Servidor escuchando en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// queryHandler maneja las consultas GET enviadas con la ruta /{año}/{mes}.
func queryHandler(w http.ResponseWriter, r *http.Request) {
	// Solo se aceptan métodos GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Se espera que la URL tenga el formato "/año/mes".
	// Se remueven los "/" extras y se separa la ruta.
	segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(segments) != 2 {
		http.Error(w, "La URL debe tener el formato /año/mes", http.StatusBadRequest)
		return
	}

	year := segments[0]
	month := segments[1]

	// Construye el patrón para la consulta SQL
	likePattern := year + month + "%" // Ejemplo: "202504%"

	// Define la consulta. Se usa @p1 para el parámetro.
	query := "SELECT Fecha, Hora, IdEmpleado FROM [ZKTime].[dbo].[FICHAJES01] WHERE Fecha LIKE @p1 ORDER BY IdEmpleado, Fecha, Hora;"

	// Ejecuta la consulta con el parámetro.
	rows, err := db.Query(query, likePattern)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error ejecutando la consulta: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Extrae los resultados y arma la lista de marcaciones.
	var results []Marcacion
	for rows.Next() {
		var m Marcacion
		// La consulta retorna las columnas en el orden: Fecha, Hora, IdEmpleado
		if err := rows.Scan(&m.Fecha, &m.Hora, &m.IdEmpleado); err != nil {
			http.Error(w, fmt.Sprintf("Error leyendo los datos: %v", err), http.StatusInternalServerError)
			return
		}
		results = append(results, m)
	}

	// Configura la cabecera de la respuesta para que sea JSON.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, fmt.Sprintf("Error codificando el JSON: %v", err), http.StatusInternalServerError)
	}
}
