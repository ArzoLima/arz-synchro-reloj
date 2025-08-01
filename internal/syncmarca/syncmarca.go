package syncmarca

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
)

// Marcacion representa cada registro extraído de la base de datos.
type Marcacion struct {
	IdEmpleado string `json:"IdEmpleado"`
	Fecha      string `json:"Fecha"`
	Hora       string `json:"Hora"`
}

// Config contiene la configuración de conexión a la base de datos
type Config struct {
	Server   string
	Database string
	UseIntegratedSecurity bool
	Username string
	Password string
	Encrypt  bool
}

// Client maneja las operaciones de sincronización de marcaciones
type Client struct {
	config Config
	db     *sql.DB
}

// NewClient crea un nuevo cliente con la configuración especificada
func NewClient(config Config) *Client {
	return &Client{
		config: config,
	}
}

// NewDefaultClient crea un cliente con la configuración por defecto
func NewDefaultClient() *Client {
	return &Client{
		config: Config{
			Server:                "SVRDP",
			Database:              "ZKTime",
			UseIntegratedSecurity: true,
			Encrypt:               false,
		},
	}
}

// Connect establece la conexión a la base de datos
func (c *Client) Connect() error {
	dsn := c.buildConnectionString()
	
	var err error
	c.db, err = sql.Open("sqlserver", dsn)
	if err != nil {
		return fmt.Errorf("error al crear la conexión: %v", err)
	}

	// Prueba la conexión con Ping.
	err = c.db.Ping()
	if err != nil {
		return fmt.Errorf("error conectando a la base de datos: %v", err)
	}

	return nil
}

// Close cierra la conexión a la base de datos
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// GetMarcaciones obtiene las marcaciones para un año y mes específicos
func (c *Client) GetMarcaciones(year, month string) ([]Marcacion, error) {
	if c.db == nil {
		return nil, fmt.Errorf("no hay conexión a la base de datos. Llama Connect() primero")
	}

	// Construye el patrón para la consulta SQL
	likePattern := year + month + "%" // Ejemplo: "202504%"

	// Define la consulta. Se usa @p1 para el parámetro.
	query := "SELECT Fecha, Hora, IdEmpleado FROM [ZKTime].[dbo].[FICHAJES01] WHERE Fecha LIKE @p1 ORDER BY IdEmpleado, Fecha, Hora;"

	// Ejecuta la consulta con el parámetro.
	rows, err := c.db.Query(query, likePattern)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando la consulta: %v", err)
	}
	defer rows.Close()

	// Extrae los resultados y arma la lista de marcaciones.
	var results []Marcacion
	for rows.Next() {
		var m Marcacion
		// La consulta retorna las columnas en el orden: Fecha, Hora, IdEmpleado
		if err := rows.Scan(&m.Fecha, &m.Hora, &m.IdEmpleado); err != nil {
			return nil, fmt.Errorf("error leyendo los datos: %v", err)
		}
		results = append(results, m)
	}

	// Verifica si hubo errores durante la iteración
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración de resultados: %v", err)
	}

	return results, nil
}

// GetMarcacionesWithAutoConnect obtiene marcaciones conectándose automáticamente
func (c *Client) GetMarcacionesWithAutoConnect(year, month string) ([]Marcacion, error) {
	dsn := fmt.Sprintf("server=%s;database=%s;integrated security=true;encrypt=disable",
		c.config.Server, c.config.Database)
	
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al crear la conexión: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error conectando a la base de datos: %v", err)
	}

	likePattern := year + month + "%"
	query := "SELECT Fecha, Hora, IdEmpleado FROM [ZKTime].[dbo].[FICHAJES01] WHERE Fecha LIKE @p1 ORDER BY IdEmpleado, Fecha, Hora;"

	rows, err := db.Query(query, likePattern)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando la consulta: %v", err)
	}
	defer rows.Close()

	var results []Marcacion
	for rows.Next() {
		var m Marcacion
		if err := rows.Scan(&m.Fecha, &m.Hora, &m.IdEmpleado); err != nil {
			return nil, fmt.Errorf("error leyendo los datos: %v", err)
		}
		results = append(results, m)
	}

	return results, rows.Err()
}

// buildConnectionString construye la cadena de conexión basada en la configuración
func (c *Client) buildConnectionString() string {
	dsn := fmt.Sprintf("server=%s;database=%s;", c.config.Server, c.config.Database)
	
	if c.config.UseIntegratedSecurity {
		dsn += "integrated security=true;"
	} else {
		dsn += fmt.Sprintf("user id=%s;password=%s;", c.config.Username, c.config.Password)
	}
	
	if c.config.Encrypt {
		dsn += "encrypt=true;"
	} else {
		dsn += "encrypt=disable;"
	}
	
	return dsn
}
