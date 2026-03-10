package main

import (
	"fmt"
	"net/http"

	"api-jugadores/controllers"
	"api-jugadores/db"

	"github.com/rs/cors"
)

func main() {
	db.Conectar()
	fs := http.FileServer(http.Dir("./uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))
	http.HandleFunc("/jugadores", controllers.ManejarJugadores)
	http.HandleFunc("/jugadores/{id}", controllers.BuscarJugador)
	http.HandleFunc("/partidos", controllers.ManejarPartidos)
	http.HandleFunc("/partidos/{id}", controllers.BuscarPartido)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})
	handler := c.Handler(http.DefaultServeMux)
	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}
