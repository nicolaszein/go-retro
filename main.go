package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nicolaszein/go-retro/handlers"
	"github.com/nicolaszein/go-retro/settings"
)

func main() {
	port := settings.PORT
	if port == "" {
		port = "8000"
	}

	db, err := gorm.Open("postgres", settings.DATABASE_URL)
	if err != nil {
		fmt.Println("Failed to connect database with error: ", err)
	}
	defer db.Close()

	fmt.Println("Starting server on port " + port)

	r := chi.NewRouter()
	r.Get("/", handlers.HealthCheck)

	env := handlers.Env{
		DB: db,
	}

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/teams", func(r chi.Router) {
			r.Post("/", env.CreateTeam)
		})
	})

	err = http.ListenAndServe(":"+port, r)

	if err != nil {
		fmt.Println("Error serving:", err)
	}
}
