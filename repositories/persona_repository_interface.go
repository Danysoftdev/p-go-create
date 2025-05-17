package repositories

import "github.com/danysoftdev/p-go-create/models"

type PersonaRepository interface {
	InsertarPersona(persona models.Persona) error
	ObtenerPersonaPorDocumento(documento string) (models.Persona, error)
}
