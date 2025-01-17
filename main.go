package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// initDB se encarga de iniciar la conexión con la base de datos
func initDB() {
	var err error
	// Usamos una IP pública o privada de tu instancia de MySQL, asegúrate de que esté accesible
	dsn := "admin:Hola1244@tcp(18.212.223.216:3306)/reservation_db"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Error al verificar conexión con la base de datos:", err)
	}
	fmt.Println("Conexión exitosa con la base de datos")
}

// Estructura para recibir datos de la reservación
type Reservation struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	// Ruta para crear una nueva reservación
	r.POST("/reservations", func(c *gin.Context) {
		var data Reservation

		// Bind JSON a la estructura de datos
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		// Validación de campos obligatorios
		if data.Name == "" || data.Email == "" || data.Date == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Todos los campos son obligatorios"})
			return
		}

		// Insertar la reservación en la base de datos
		_, err := db.Exec("INSERT INTO reservations (name, email, reservation_date) VALUES (?, ?, ?)",
			data.Name, data.Email, data.Date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la reservación"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Reservación confirmada"})
	})

	// Iniciar el servidor en el puerto 8080
	r.Run(":8080")
}
