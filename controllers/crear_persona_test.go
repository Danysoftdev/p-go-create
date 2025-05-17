package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danysoftdev/p-go-create/controllers"
	"github.com/danysoftdev/p-go-create/models"
	"github.com/danysoftdev/p-go-create/services"
	"github.com/danysoftdev/p-go-create/tests/mocks"
	"github.com/stretchr/testify/assert"
)


func TestCrearPersonaController_Success(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	persona := models.Persona{
		Documento: "123",
		Nombre:    "Laura",
		Apellido:  "Gómez",
		Edad:      30,
		Correo:    "laura@example.com",
		Telefono:  "123456789",
		Direccion: "Calle Falsa 123",
	}

	// Mock de flujo exitoso: no existe, y se inserta correctamente
	mockRepo.On("ObtenerPersonaPorDocumento", "123").Return(models.Persona{}, errors.New("no encontrado"))
	mockRepo.On("InsertarPersona", persona).Return(nil)

	body, _ := json.Marshal(persona)
	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	controllers.CrearPersona(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "Persona creada exitosamente")
	mockRepo.AssertExpectations(t)
}

func TestCrearPersonaController_DocumentoExistente(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	existente := models.Persona{
		Documento: "123",
		Nombre:    "Existente",
		Apellido:  "Persona",
		Edad:      25,
		Correo:    "existente@example.com",
		Telefono:  "987654321",
		Direccion: "Otra calle",
	}

	mockRepo.On("ObtenerPersonaPorDocumento", "123").Return(existente, nil)

	body, _ := json.Marshal(existente)
	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	controllers.CrearPersona(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "ya existe una persona con ese documento")
	mockRepo.AssertExpectations(t)
}

func TestCrearPersonaController_JSONInvalido(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBuffer([]byte("{invalido")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	controllers.CrearPersona(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "El formato del cuerpo es inválido")
}