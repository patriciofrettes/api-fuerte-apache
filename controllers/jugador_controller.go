package controllers

import (
	"api-jugadores/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

var jugadores = []models.Jugador{
	{
		ID:       1,
		Nombre:   "Patricio Frettes",
		Posicion: "Volante",
		Edad:     24,
	},
	{
		ID:       2,
		Nombre:   "Rodrigo Cordoba",
		Posicion: "Defensor",
		Edad:     25,
	},
}

func ManejarJugadores(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jugadores)
	} else if r.Method == "POST" {
		var nuevoJugador models.Jugador
		json.NewDecoder(r.Body).Decode(&nuevoJugador)
		jugadores = append(jugadores, nuevoJugador)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Jugador agregado al plantel")
	}
}

func BuscarJugador(w http.ResponseWriter, r *http.Request) {
	idTexto := r.PathValue("id")
	idNumero, err := strconv.Atoi(idTexto)
	if err != nil {
		http.Error(w, "El ID debe ser un número", http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		for _, jugador := range jugadores {
			if jugador.ID == idNumero {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(jugador)
				return
			}
		}
		http.Error(w, "Jugador no encontrado en el plantel", http.StatusNotFound)
	} else if r.Method == "PUT" {
		var jugadorActualizado models.Jugador
		json.NewDecoder(r.Body).Decode(&jugadorActualizado)
		for i, jugador := range jugadores {
			if jugador.ID == idNumero {
				jugadorActualizado.ID = jugador.ID
				jugadores[i] = jugadorActualizado
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "Jugador actualizado con éxito")
				return
			}
		}
		http.Error(w, "Jugador no encontrado para actualizar", http.StatusNotFound)
	} else if r.Method == "DELETE" {
		for i, jugador := range jugadores {
			if jugador.ID == idNumero {
				jugadores = append(jugadores[:i], jugadores[i+1:]...)
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "Jugador eliminado del plantel")
				return
			}
		}
		http.Error(w, "Jugador no encontrado para eliminar", http.StatusNotFound)
	}
}
