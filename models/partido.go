package models

type Partido struct {
	ID        int    `json:"id"`
	Rival     string `json:"rival"`
	Fecha     string `json:"fecha"`
	Resultado string `json:"resultado"`
}
