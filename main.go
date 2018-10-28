package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
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
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	r.Get("/", handlers.HealthCheck)

	env := handlers.Env{
		DB: db,
	}

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/teams", func(r chi.Router) {
			r.Get("/", env.ListTeams)
			r.Get("/{teamID}", env.FetchTeam)
			r.Post("/", env.CreateTeam)
		})

		r.Route("/retrospectives", func(r chi.Router) {
			r.Post("/", env.CreateRetrospective)
			r.Get("/{retrospectiveID}", env.FetchRetrospective)

			// Cards
			r.Route("/{retrospectiveID}/cards", func(r chi.Router) {
				r.Post("/", env.CreateCard)
				r.Post("/{cardID}/votes", env.AddCardVote)
			})
		})
	})

	err = http.ListenAndServe(":"+port, r)

	if err != nil {
		fmt.Println("Error serving:", err)
	}
}
