package models

type Jugador struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Posicion string `json:"posicion"`
	Edad     int    `json:"edad"`
	Foto     string `json:"foto"`
}
