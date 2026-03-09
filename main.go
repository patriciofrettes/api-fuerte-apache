package main

import (
	"fmt"
	"net/http"

	"api-jugadores/controllers"
	"api-jugadores/db"
)

func main() {
	db.Conectar()
	http.HandleFunc("/jugadores", controllers.ManejarJugadores)
	http.HandleFunc("/jugadores/{id}", controllers.BuscarJugador)
	http.HandleFunc("/partidos", controllers.ManejarPartidos)
	http.HandleFunc("/partidos/{id}", controllers.BuscarPartido)
	fmt.Println("Servidor de Fuerte Apache corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
