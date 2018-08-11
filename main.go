package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nicolaszein/go-retro/database"
	"github.com/nicolaszein/go-retro/handlers"
	"github.com/nicolaszein/go-retro/settings"
)

func main() {
	port := settings.PORT
	if port == "" {
		port = "8000"
	}

	db, err := database.NewPostgres(settings.DATABASE_URL)
	if err != nil {
		fmt.Println("Failed to connect database with error: ", err)
	}
	defer db.DB.Close()

	fmt.Println("Starting server on port " + port)

	r := chi.NewRouter()
	r.Get("/", handlers.HealthCheck)

	env := handlers.Env{
		DB: db,
	}

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/teams", func(r chi.Router) {
			r.Get("/", env.ListTeams)
			r.Post("/", env.CreateTeam)
		})

		r.Route("/retrospectives", func(r chi.Router) {
			r.Post("/", env.CreateRetrospective)
			r.Post("/{retrospectiveID}/cards", env.CreateCard)
		})
	})

	err = http.ListenAndServe(":"+port, r)

	if err != nil {
		fmt.Println("Error serving:", err)
	}
}
