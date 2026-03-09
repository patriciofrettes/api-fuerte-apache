package controllers

import (
	"api-jugadores/db"
	"api-jugadores/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func ManejarJugadores(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		filas, err := db.DB.Query("SELECT id, nombre, posicion, edad FROM jugadores")
		if err != nil {
			http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
			return
		}
		defer filas.Close()

		var jugadores []models.Jugador

		for filas.Next() {
			var jugador models.Jugador
			filas.Scan(&jugador.ID, &jugador.Nombre, &jugador.Posicion, &jugador.Edad)
			jugadores = append(jugadores, jugador)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jugadores)
	} else if r.Method == "POST" {
		var nuevoJugador models.Jugador
		json.NewDecoder(r.Body).Decode(&nuevoJugador)
		_, err := db.DB.Exec("INSERT INTO jugadores (nombre, posicion, edad) VALUES (?, ?, ?)", nuevoJugador.Nombre, nuevoJugador.Posicion, nuevoJugador.Edad)
		if err != nil {
			http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "¡Jugador agregado al plantel de Fuerte Apache en MySQL!")
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
		var jugador models.Jugador
		err := db.DB.QueryRow("SELECT id, nombre, posicion, edad FROM jugadores WHERE id = ?", idNumero).Scan(&jugador.ID, &jugador.Nombre, &jugador.Posicion, &jugador.Edad)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Jugador no encontrado en el plantel", http.StatusNotFound)
			} else {
				http.Error(w, "Error al buscar en la base de datos", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jugador)

	} else if r.Method == "PUT" {
		var jugadorActualizado models.Jugador
		json.NewDecoder(r.Body).Decode(&jugadorActualizado)

		_, err := db.DB.Exec("UPDATE jugadores SET nombre = ?, posicion = ?, edad = ? WHERE id = ?",
			jugadorActualizado.Nombre, jugadorActualizado.Posicion, jugadorActualizado.Edad, idNumero)

		if err != nil {
			http.Error(w, "Error al actualizar en la base de datos", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Jugador actualizado con éxito")

	} else if r.Method == "DELETE" {
		_, err := db.DB.Exec("DELETE FROM jugadores WHERE id = ?", idNumero)
		if err != nil {
			http.Error(w, "Error al eliminar en la base de datos", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Jugador eliminado del plantel")
	}
}
