package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jacob-tech-challenge/api/handlers"
	// Import the handler functions from the course package
)

// SetupRoutes sets up the API routes using the Chi router.
func SetupRoutes(db *sql.DB) (http.Handler) {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Mount("/course", courseRoutes(db));
		// r.Mount("/person", personRoutes(db));
	})

	return r
}

// courseRoutes defines the routes for the /api/course endpoint.
func courseRoutes(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Get("/", handlers.HandleGetAllCourses(db))
	r.Get("/{id}", handlers.HandleGetCourseByID(db))
	r.Put("/{id}", handlers.HandleUpdateCourse(db))
	// r.Post("/", createCourse)
	// r.Delete("/{id}", deleteCourseByID)

	return r
}

// personRoutes defines the routes for the /api/person endpoint.
// func personRoutes() http.Handler {
// 	r := chi.NewRouter()

// 	r.Get("/", getAllPersons) // Note: This route will handle query parameters
// 	r.Get("/{name}", getPersonByName)
// 	r.Put("/{name}", updatePersonByName)
// 	r.Post("/", createPerson)
// 	r.Delete("/{name}", deletePersonByName)

// 	return r
// }

// ... Implement your handler functions for each route here (getAllCourses, getCourseByID, etc.)