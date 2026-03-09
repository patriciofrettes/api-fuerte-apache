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

func ManejarPartidos(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		filas, err := db.DB.Query("SELECT id, rival, fecha, resultado FROM partidos")
		if err != nil {
			http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
			return
		}
		defer filas.Close()

		var partidos []models.Partido

		for filas.Next() {
			var partido models.Partido
			filas.Scan(&partido.ID, &partido.Rival, &partido.Fecha, &partido.Resultado)
			partidos = append(partidos, partido)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(partidos)

	} else if r.Method == "POST" {
		var nuevoPartido models.Partido
		json.NewDecoder(r.Body).Decode(&nuevoPartido)

		_, err := db.DB.Exec("INSERT INTO partidos (rival, fecha, resultado) VALUES (?, ?, ?)", nuevoPartido.Rival, nuevoPartido.Fecha, nuevoPartido.Resultado)
		if err != nil {
			http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "¡Partido contra %s agregado al fixture de Fuerte Apache!", nuevoPartido.Rival)
	}
}

func BuscarPartido(w http.ResponseWriter, r *http.Request) {
	idTexto := r.PathValue("id")
	idNumero, err := strconv.Atoi(idTexto)
	if err != nil {
		http.Error(w, "El ID debe ser un número", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		var partido models.Partido
		err := db.DB.QueryRow("SELECT id, rival, fecha, resultado FROM partidos WHERE id = ?", idNumero).Scan(&partido.ID, &partido.Rival, &partido.Fecha, &partido.Resultado)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Partido no encontrado", http.StatusNotFound)
			} else {
				http.Error(w, "Error al buscar en la base de datos", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(partido)

	} else if r.Method == "PUT" {
		var partidoActualizado models.Partido
		json.NewDecoder(r.Body).Decode(&partidoActualizado)

		_, err := db.DB.Exec("UPDATE partidos SET rival = ?, fecha = ?, resultado = ? WHERE id = ?",
			partidoActualizado.Rival, partidoActualizado.Fecha, partidoActualizado.Resultado, idNumero)

		if err != nil {
			http.Error(w, "Error al actualizar el partido", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Partido actualizado con éxito")

	} else if r.Method == "DELETE" {
		_, err := db.DB.Exec("DELETE FROM partidos WHERE id = ?", idNumero)
		if err != nil {
			http.Error(w, "Error al eliminar el partido", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Partido eliminado del fixture")
	}
}
