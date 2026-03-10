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
		filas, err := db.DB.Query("SELECT id, rival, fecha, resultado FROM partidos ORDER BY fecha DESC")
		if err != nil {
			http.Error(w, "Error al consultar partidos", http.StatusInternalServerError)
			return
		}
		defer filas.Close()

		var partidos []models.Partido
		for filas.Next() {
			var p models.Partido
			filas.Scan(&p.ID, &p.Rival, &p.Fecha, &p.Resultado)
			partidos = append(partidos, p)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(partidos)

	} else if r.Method == "POST" {
		var nuevo models.Partido
		err := json.NewDecoder(r.Body).Decode(&nuevo)
		if err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		_, err = db.DB.Exec("INSERT INTO partidos (rival, fecha, resultado) VALUES (?, ?, ?)",
			nuevo.Rival, nuevo.Fecha, nuevo.Resultado)

		if err != nil {
			http.Error(w, "Error al guardar el partido", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Partido programado con éxito")
	}
}

func BuscarPartido(w http.ResponseWriter, r *http.Request) {
	idTexto := r.PathValue("id")
	idNumero, _ := strconv.Atoi(idTexto)

	if r.Method == "GET" {
		var p models.Partido
		err := db.DB.QueryRow("SELECT id, rival, fecha, resultado FROM partidos WHERE id = ?", idNumero).
			Scan(&p.ID, &p.Rival, &p.Fecha, &p.Resultado)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Partido no encontrado", http.StatusNotFound)
			} else {
				http.Error(w, "Error en la base de datos", http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)

	} else if r.Method == "PUT" {
		var p models.Partido
		json.NewDecoder(r.Body).Decode(&p)

		_, err := db.DB.Exec("UPDATE partidos SET rival = ?, fecha = ?, resultado = ? WHERE id = ?",
			p.Rival, p.Fecha, p.Resultado, idNumero)

		if err != nil {
			http.Error(w, "Error al actualizar partido", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Partido actualizado")

	} else if r.Method == "DELETE" {
		_, err := db.DB.Exec("DELETE FROM partidos WHERE id = ?", idNumero)
		if err != nil {
			http.Error(w, "Error al eliminar partido", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Partido borrado del fixture")
	}
}
