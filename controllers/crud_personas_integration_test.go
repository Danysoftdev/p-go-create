//go:build integration
// +build integration

package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/danysoftdev/p-go-create/config"
	"github.com/danysoftdev/p-go-create/controllers"
	"github.com/danysoftdev/p-go-create/models"
	"github.com/danysoftdev/p-go-create/repositories"
	"github.com/danysoftdev/p-go-create/services"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestEndpointsControllerIntegration(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(20 * time.Second),
	}
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	defer mongoC.Terminate(ctx)

	endpoint, err := mongoC.Endpoint(ctx, "")
	assert.NoError(t, err)

	os.Setenv("MONGO_URI", "mongodb://"+endpoint)
	os.Setenv("MONGO_DB", "testdb")
	os.Setenv("COLLECTION_NAME", "personas_test")

	err = config.ConectarMongo()
	assert.NoError(t, err)
	defer config.CerrarMongo()

	repositories.SetCollection(config.Collection)
	services.SetPersonaRepository(repositories.RealPersonaRepository{})

	// Setup router
	router := mux.NewRouter()
	router.HandleFunc("/personas", controllers.CrearPersona).Methods("POST")
	
	// 1. Crear persona
	persona := models.Persona{
		Documento: "999",
		Nombre:    "Test",
		Apellido:  "Integration",
		Edad:      33,
		Correo:    "test@integration.com",
		Telefono:  "1111111111",
		Direccion: "Calle Test",
	}
	body, _ := json.Marshal(persona)
	reqCrear := httptest.NewRequest("POST", "/personas", bytes.NewReader(body))
	resCrear := httptest.NewRecorder()
	router.ServeHTTP(resCrear, reqCrear)

	assert.Equal(t, http.StatusCreated, resCrear.Code)
}
