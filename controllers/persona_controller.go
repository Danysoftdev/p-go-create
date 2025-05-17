package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/danysoftdev/p-go-create/models"
	"github.com/danysoftdev/p-go-create/services"
)

func CrearPersona(w http.ResponseWriter, r *http.Request) {
	var persona models.Persona

	err := json.NewDecoder(r.Body).Decode(&persona)
	if err != nil {
		http.Error(w, "El formato del cuerpo es inválido", http.StatusBadRequest)
		return
	}

	err = services.CrearPersona(persona)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Persona creada exitosamente"})
}