package repositories

import (
	"context"
	"time"

	"github.com/danysoftdev/p-go-create/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

// Permite inyectar la colección desde fuera (ideal para pruebas)
func SetCollection(c *mongo.Collection) {
	collection = c
}

// InsertarPersona guarda una nueva persona en la base de datos
func InsertarPersona(persona models.Persona) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, persona)
	return err
}

// ObtenerPersonaPorDocumento busca una persona por su Documento
func ObtenerPersonaPorDocumento(documento string) (models.Persona, error) {
	var persona models.Persona
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"documento": documento}).Decode(&persona)
	return persona, err
}

type RealPersonaRepository struct{}

func (r RealPersonaRepository) InsertarPersona(p models.Persona) error {
	return InsertarPersona(p)
}

func (r RealPersonaRepository) ObtenerPersonaPorDocumento(doc string) (models.Persona, error) {
	return ObtenerPersonaPorDocumento(doc)
}