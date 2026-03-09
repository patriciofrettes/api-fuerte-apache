package main

import (
	"fmt"
	"net/http"

	"api-jugadores/controllers"
)

func main() {
	http.HandleFunc("/jugadores", controllers.ManejarJugadores)
	http.HandleFunc("/jugadores/{id}", controllers.BuscarJugador)

	fmt.Println("Servidor de Fuerte Apache corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
