package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
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
		if _, err := services.UpdateCourse(db, id, course); err != nil {
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

func HandleCreateCourse(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var course models.Course
		if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, err := services.CreateCourse(db, course); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		courseOut := map[string]interface{}{
			"id":   course.ID,
			"name": course.Name,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(courseOut); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func HandleDeleteCourse(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := services.DeleteCourse(db, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		// json formatted string representing deletion confirmation
		w.Write([]byte(`{"message": "Course deleted"}`))
	})
}

func HandleGetAllPeople(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// chi parameters, if any
		name := r.URL.Query().Get("name")
		age := 0
		ageString := r.URL.Query().Get("age")

		log.Println("name:", name)
		log.Println("ageString:", ageString)


		if ageString != "" {
			var err error
			age, err = strconv.Atoi(ageString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}	

		people, err := services.GetAllPeople(db, name, age)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		peopleOut := make([]map[string]interface{}, len(people))
		for i, person := range people {
			peopleOut[i] = map[string]interface{}{
				"id":        person.ID,
				"firstName": person.FirstName,
				"lastName":  person.LastName,
				"type":      person.Type,
				"age":       person.Age,
				"courses":   person.Courses,
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(peopleOut); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func HandleGetPersonByName(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		person, err := services.GetPersonByName(db, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		personOut := map[string]interface{}{
			"id":        person.ID,
			"firstName": person.FirstName,
			"lastName":  person.LastName,
			"type":      person.Type,
			"age":       person.Age,
			"courses":   person.Courses,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(personOut); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func HandleUpdatePersonByName(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		if name == "" {
			http.Error(w, "name parameter is required", http.StatusBadRequest)
			return
		}

		var person models.Person
		if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Check if person exists first
		var existingPerson models.Person
		err := db.QueryRow("SELECT id, first_name, last_name, type, age FROM person WHERE first_name = $1", name).
			Scan(&existingPerson.ID, &existingPerson.FirstName, &existingPerson.LastName, &existingPerson.Type, &existingPerson.Age)
		if err == sql.ErrNoRows {
			http.Error(w, "Person not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Update person
		_, err = db.Exec(`
			UPDATE person 
			SET first_name = $1, last_name = $2, type = $3, age = $4 
			WHERE first_name = $5`,
			person.FirstName, person.LastName, person.Type, person.Age, name)
		if err != nil {
			http.Error(w, "Failed to update person: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Handle courses
		if len(person.Courses) > 0 {
			for _, courseID := range person.Courses {
				_, err = db.Exec(`
					INSERT INTO person_course (person_id, course_id) 
					VALUES ($1, $2)
					ON CONFLICT (person_id, course_id) DO NOTHING`,
					existingPerson.ID, courseID)
				if err != nil {
					http.Error(w, "Failed to update courses: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		// Prepare response
		person.ID = existingPerson.ID // Make sure we include the correct ID
		personOutMap := map[string]interface{}{
			"id":        float64(person.ID),
			"firstName": person.FirstName,
			"lastName":  person.LastName,
			"type":      person.Type,
			"age":       float64(person.Age),
			"courses":   person.Courses,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(personOutMap); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	})
}

func HandleCreatePerson(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var person models.Person
		if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validate person
		if person.FirstName == "" || person.LastName == "" {
			http.Error(w, "FirstName and LastName are required", http.StatusBadRequest)
			return
		}

		// Create person
		personOut, err := services.CreatePerson(db, person)
		if err != nil {
			log.Printf("Error creating person: %v", err) // Add logging
			http.Error(w, "Failed to create person: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Update personOut with the courses from the input
		personOut.Courses = person.Courses

		personOutMap := map[string]interface{}{
			"id":        personOut.ID,
			"firstName": personOut.FirstName,
			"lastName":  personOut.LastName,
			"type":      personOut.Type,
			"age":       personOut.Age,
			"courses":   personOut.Courses,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(personOutMap); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
}

func HandleDeletePersonByName(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		if err := services.DeletePersonByName(db, name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		// json formatted string representing deletion confirmation
		w.Write([]byte(`{"message": "Person deleted"}`))
	})
}


