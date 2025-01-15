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

func initDB() {
    var err error
    dsn := "admin:securepassword@tcp(<DB_EC2_IP>:3306)/reservation_db"
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Error al conectar a la base de datos:", err)
    }
    if err = db.Ping(); err != nil {
        log.Fatal("Error al verificar conexión con la base de datos:", err)
    }
    fmt.Println("Conexión exitosa con la base de datos")
}

func main() {
    initDB()
    defer db.Close()

    r := gin.Default()

    r.POST("/reservations", func(c *gin.Context) {
        var data struct {
            Name  string `json:"name"`
            Email string `json:"email"`
            Date  string `json:"date"`
        }
        if err := c.ShouldBindJSON(&data); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
            return
        }

        _, err := db.Exec("INSERT INTO reservations (name, email, reservation_date) VALUES (?, ?, ?)",
            data.Name, data.Email, data.Date)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la reservación"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "Reservación confirmada"})
    })

    r.Run(":8080")
}
