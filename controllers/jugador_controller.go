package controllers

import (
	"api-jugadores/db"
	"api-jugadores/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func ManejarJugadores(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Seleccionamos también la columna foto
		filas, err := db.DB.Query("SELECT id, nombre, posicion, edad, foto FROM jugadores")
		if err != nil {
			http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
			return
		}
		defer filas.Close()

		var jugadores []models.Jugador

		for filas.Next() {
			var jugador models.Jugador
			// Escaneamos el resultado incluyendo la foto
			errScan := filas.Scan(&jugador.ID, &jugador.Nombre, &jugador.Posicion, &jugador.Edad, &jugador.Foto)
			if errScan != nil {
				continue
			}
			jugadores = append(jugadores, jugador)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jugadores)

	} else if r.Method == "POST" {
		r.ParseMultipartForm(10 << 20)
		nombre := r.FormValue("nombre")
		posicion := r.FormValue("posicion")
		edad, _ := strconv.Atoi(r.FormValue("edad"))

		file, handler, err := r.FormFile("foto")
		var rutaFoto string
		if err == nil {
			defer file.Close()
			rutaGuardado := "uploads/" + handler.Filename
			archivoDestino, errArchivo := os.Create(rutaGuardado)

			if errArchivo == nil {
				defer archivoDestino.Close()
				io.Copy(archivoDestino, file)
				rutaFoto = "http://localhost:8080/uploads/" + handler.Filename
			}
		}

		_, errDB := db.DB.Exec("INSERT INTO jugadores (nombre, posicion, edad, foto) VALUES (?, ?, ?, ?)", nombre, posicion, edad, rutaFoto)
		if errDB != nil {
			http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "¡Jugador agregado con su foto al plantel!")
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
		err := db.DB.QueryRow("SELECT id, nombre, posicion, edad, foto FROM jugadores WHERE id = ?", idNumero).
			Scan(&jugador.ID, &jugador.Nombre, &jugador.Posicion, &jugador.Edad, &jugador.Foto)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Jugador no encontrado", http.StatusNotFound)
			} else {
				http.Error(w, "Error en la base de datos", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jugador)

	} else if r.Method == "PUT" {
		// --- ACTUALIZACIÓN PARA SOPORTAR FOTOS AL EDITAR ---
		r.ParseMultipartForm(10 << 20)
		nombre := r.FormValue("nombre")
		posicion := r.FormValue("posicion")
		edad, _ := strconv.Atoi(r.FormValue("edad"))

		file, handler, errFile := r.FormFile("foto")

		var errQuery error
		if errFile == nil {
			// Si se subió una foto nueva, la guardamos y actualizamos todo
			defer file.Close()
			rutaGuardado := "uploads/" + handler.Filename
			archivoDestino, _ := os.Create(rutaGuardado)
			defer archivoDestino.Close()
			io.Copy(archivoDestino, file)
			nuevaRuta := "http://localhost:8080/uploads/" + handler.Filename

			_, errQuery = db.DB.Exec("UPDATE jugadores SET nombre = ?, posicion = ?, edad = ?, foto = ? WHERE id = ?",
				nombre, posicion, edad, nuevaRuta, idNumero)
		} else {
			// Si NO se subió foto, actualizamos solo los datos de texto
			_, errQuery = db.DB.Exec("UPDATE jugadores SET nombre = ?, posicion = ?, edad = ? WHERE id = ?",
				nombre, posicion, edad, idNumero)
		}

		if errQuery != nil {
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
