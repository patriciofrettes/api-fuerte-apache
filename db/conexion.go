package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Conectar() {
	dsn := "root:@tcp(127.0.0.1:3306)/fuerte_apache"
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic("Error al abrir la base de datos: " + err.Error())
	}
	err = DB.Ping()
	if err != nil {
		panic("Error al conectar a la base de datos (¿Está prendido XAMPP?): " + err.Error())
	}
	fmt.Println("¡Conexión exitosa a la base de datos de Fuerte Apache! ⚽")
}
