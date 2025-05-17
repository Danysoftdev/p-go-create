package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/danysoftdev/p-go-create/config"
	"github.com/danysoftdev/p-go-create/controllers"
	"github.com/danysoftdev/p-go-create/repositories"
	"github.com/danysoftdev/p-go-create/services"

	"github.com/gorilla/mux"
)

func main() {
	// Conectamos a MongoDB
	err := config.ConectarMongo()
	if err != nil {
		log.Fatal("❌ Error conectando a MongoDB:", err)
	}


	// 2. Inyectar el repositorio real
	services.SetPersonaRepository(repositories.RealPersonaRepository{})

	// 3. Inyectar la colección de MongoDB
	repositories.SetCollection(config.Collection)

	// Creamos el enrutador
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hola, desde la creación de personas")
	})

	// Rutas de la API
	router.HandleFunc("/crear-personas", controllers.CrearPersona).Methods("POST")
	
	// Puerto de escucha
	puerto := ":8080"
	fmt.Printf("🚀 Servidor escuchando en http://localhost%s\n", puerto)
	log.Fatal(http.ListenAndServe(puerto, router))
}
