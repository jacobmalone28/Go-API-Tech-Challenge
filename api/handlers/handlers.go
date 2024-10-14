package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jacob-tech-challenge/api/models"
	"github.com/jacob-tech-challenge/api/services"
)

func HandleGetAllCourses(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		courses, err := services.GetAllCourses(db)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		coursesOut := make([]map[string]interface{}, len(courses))
		for i, course := range courses {
			coursesOut[i] = map[string]interface{}{
				"id":   course.ID,
				"name": course.Name,
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(coursesOut); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// HandleGetCourseByID handles the get course by id request
func HandleGetCourseByID(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		course, err := services.GetCourseByID(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		courseOut := map[string]interface{}{
			"id":   course.ID,
			"name": course.Name,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(courseOut); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	
}

func HandleUpdateCourse(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var course models.Course
		if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		course.ID = id
		if _, err := services.UpdateCourse(db, course); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		courseOut := map[string]interface{}{
			"id":   course.ID,
			"name": course.Name,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(courseOut); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}