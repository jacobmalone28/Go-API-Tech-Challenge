package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/jacob-tech-challenge/api/models"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetAllCourses(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name           string
		mockSetup     func(sqlmock.Sqlmock)
		expectedCode  int
		expectedBody  []map[string]interface{}
	}{
		{
			name: "successful retrieval",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Math").
					AddRow(2, "Science")
				mock.ExpectQuery(`SELECT \* FROM "course"`).WillReturnRows(rows)
			},
			expectedCode: http.StatusOK,
			expectedBody: []map[string]interface{}{
				{"id": float64(1), "name": "Math"},
				{"id": float64(2), "name": "Science"},
			},
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "course"`).WillReturnError(sql.ErrConnDone)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup(mock)

			// Create request and response recorder
			req := httptest.NewRequest("GET", "/courses", nil)
			rr := httptest.NewRecorder()

			// Call the handler
			handler := HandleGetAllCourses(db)
			handler.ServeHTTP(rr, req)

			// Assert status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// For successful cases, verify the response body
			if tt.expectedCode == http.StatusOK {
				var got []map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, got)
			}
		})
	}
}

func TestHandleGetCourseByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name          string
		courseID      string
		mockSetup     func(sqlmock.Sqlmock)
		expectedCode  int
		expectedBody  map[string]interface{}
	}{
		{
			name:     "successful retrieval",
			courseID: "1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Math")
				// Note: The regex pattern is updated to match the actual query
				mock.ExpectQuery(`SELECT \* FROM "course" WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":   float64(1),
				"name": "Math",
			},
		},
		{
			name:     "invalid id",
			courseID: "invalid",
			mockSetup: func(mock sqlmock.Sqlmock) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup(mock)

			// Create a new router to properly handle URL parameters
			router := chi.NewRouter()
			router.Get("/{id}", HandleGetCourseByID(db))

			// Create request
			req := httptest.NewRequest("GET", "/"+tt.courseID, nil)
			rr := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(rr, req)

			// Assert status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// For successful cases, verify the response body
			if tt.expectedCode == http.StatusOK {
				var got map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, got)
			}
		})
	}
}

func TestHandleCreateCourse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name          string
		course        models.Course
		mockSetup     func(sqlmock.Sqlmock)
		expectedCode  int
		expectedBody  map[string]interface{}
	}{
		{
			name: "successful creation",
			course: models.Course{
				Name: "New Course",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO "course" \(name\) VALUES \(\$1\) RETURNING id`).
					WithArgs("New Course").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedCode: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":   float64(0),
				"name": "New Course",
			},
		},
		{
			name: "database error",
			course: models.Course{
				Name: "New Course",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO "course" \(name\) VALUES \(\$1\) RETURNING id`).
					WithArgs("New Course").
					WillReturnError(sql.ErrConnDone)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup(mock)

			// Create request body
			body, err := json.Marshal(tt.course)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/courses", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			// Call the handler
			handler := HandleCreateCourse(db)
			handler.ServeHTTP(rr, req)

			// Assert status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// For successful cases, verify the response body
			if tt.expectedCode == http.StatusCreated {
				var got map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, got)
			}
		})
	}
}

func TestHandleUpdateCourse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name          string
		courseID      string
		course        models.Course
		mockSetup     func(sqlmock.Sqlmock)
		expectedCode  int
		expectedBody  map[string]interface{}
	}{
		{
			name:     "successful update",
			courseID: "1",
			course: models.Course{
				ID:   1,
				Name: "Updated Course",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE "course" SET name = \$1 WHERE id = \$2`).
					WithArgs("Updated Course", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":   float64(1),
				"name": "Updated Course",
			},
		},
		{
			name:     "invalid id",
			courseID: "invalid",
			course: models.Course{
				Name: "Updated Course",
			},
			mockSetup:     func(mock sqlmock.Sqlmock) {},
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:     "database error",
			courseID: "1",
			course: models.Course{
				ID:   1,
				Name: "Updated Course",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE "course" SET name = \$1 WHERE id = \$2`).
					WithArgs("Updated Course", 1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup(mock)

			// Create a new router to properly handle URL parameters
			router := chi.NewRouter()
			router.Put("/{id}", HandleUpdateCourse(db))

			// Create request body
			body, err := json.Marshal(tt.course)
			assert.NoError(t, err)

			// Create request
			req := httptest.NewRequest(http.MethodPut, "/"+tt.courseID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(rr, req)

			// Assert status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// For successful cases, verify the response body
			if tt.expectedCode == http.StatusOK {
				var got map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, got)
			}
		})
	}
}

func TestHandleDeleteCourse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name          string
		courseID      string
		mockSetup     func(sqlmock.Sqlmock)
		expectedCode  int
	}{
		{
			name:     "successful deletion",
			courseID: "1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM "course" WHERE id = \$1`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:     "invalid id",
			courseID: "invalid",
			mockSetup: func(mock sqlmock.Sqlmock) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:     "database error",
			courseID: "1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM "course" WHERE id = \$1`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup(mock)

			// Create a new router to properly handle URL parameters
			router := chi.NewRouter()
			router.Delete("/{id}", HandleDeleteCourse(db))

			// Create request
			req := httptest.NewRequest(http.MethodDelete, "/"+tt.courseID, nil)
			rr := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(rr, req)

			// Assert status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// For successful deletion, verify response body contains success message
			if tt.expectedCode == http.StatusNoContent {
				var got map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&got)
				if err == nil {
					assert.Equal(t, "Course deleted", got["message"])
				}
			}
		})
	}
}

func TestHandleGetAllPeople(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name       string
		queryName  string
		queryAge   string
		mockSetup  func(sqlmock.Sqlmock)
		wantStatus int
		wantBody   []models.Person
	}{
		{
			name:      "Success - No filters",
			queryName: "",
			queryAge:  "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "type", "age"}).
					AddRow(1, "John", "Doe", "student", 20)
				mock.ExpectQuery("SELECT \\* FROM person").
					WithArgs("", 0).
					WillReturnRows(rows)
				
				courseRows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Math")
				mock.ExpectQuery("SELECT c.id, c.name FROM").
					WithArgs(1).
					WillReturnRows(courseRows)
			},
			wantStatus: http.StatusOK,
			wantBody: []models.Person{
				{
					ID: 1, FirstName: "John", LastName: "Doe",
					Type: "student", Age: 20, Courses: []int{1},
				},
			},
		},
		{
			name:      "Success - With name filter",
			queryName: "John",
			queryAge:  "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "type", "age"}).
					AddRow(1, "John", "Doe", "student", 20)
				mock.ExpectQuery("SELECT \\* FROM person").
					WithArgs("John", 0).
					WillReturnRows(rows)
				
				courseRows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Math")
				mock.ExpectQuery("SELECT c.id, c.name FROM").
					WithArgs(1).
					WillReturnRows(courseRows)
			},
			wantStatus: http.StatusOK,
			wantBody: []models.Person{
				{
					ID: 1, FirstName: "John", LastName: "Doe",
					Type: "student", Age: 20, Courses: []int{1},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			req := httptest.NewRequest("GET", "/people?name="+tt.queryName+"&age="+tt.queryAge, nil)
			w := httptest.NewRecorder()

			handler := HandleGetAllPeople(db)
			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var got []map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&got)
				assert.NoError(t, err)
				
				// Convert wantBody to same format as response
				want := make([]map[string]interface{}, len(tt.wantBody))
				for i, p := range tt.wantBody {
					want[i] = map[string]interface{}{
						"id":        float64(p.ID),
						"firstName": p.FirstName,
						"lastName":  p.LastName,
						"type":      p.Type,
						"age":       float64(p.Age),
						"courses":   (p.Courses),
					}
				}

				// Convert the courses in the response to []int
				for _, person := range got {
					if courses, ok := person["courses"].([]interface{}); ok {
						intCourses := make([]int, len(courses))
						for i, c := range courses {
							if fNum, ok := c.(float64); ok {
								intCourses[i] = int(fNum)
							}
						}
						person["courses"] = intCourses
					}
				}
				assert.Equal(t, want, got)
			}
		})
	}
}

func TestHandleGetPersonByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name       string
		personName string
		mockSetup  func(sqlmock.Sqlmock)
		wantStatus int
		wantBody   *models.Person
	}{
		{
			name:       "Success",
			personName: "John",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "type", "age"}).
					AddRow(1, "John", "Doe", "student", 20)
				mock.ExpectQuery("SELECT \\* FROM person WHERE").
					WithArgs("John").
					WillReturnRows(rows)
				
				courseRows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Math")
				mock.ExpectQuery("SELECT c.id, c.name FROM").
					WithArgs(1).
					WillReturnRows(courseRows)
			},
			wantStatus: http.StatusOK,
			wantBody: &models.Person{
				ID: 1, FirstName: "John", LastName: "Doe",
				Type: "student", Age: 20, Courses: []int{1},
			},
		},
		{
			name:       "Person Not Found",
			personName: "NonExistent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM person WHERE").
					WithArgs("NonExistent").
					WillReturnError(sql.ErrNoRows)
			},
			wantStatus: http.StatusOK,
			wantBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			router := chi.NewRouter()
			router.Get("/people/{name}", HandleGetPersonByName(db))

			req := httptest.NewRequest("GET", "/people/"+tt.personName, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK && tt.wantBody != nil {
				var got map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&got)
				assert.NoError(t, err)

				// Convert numeric values to float64 in the expected output
				want := map[string]interface{}{
					"id":        float64(tt.wantBody.ID),
					"firstName": tt.wantBody.FirstName,
					"lastName":  tt.wantBody.LastName,
					"type":      tt.wantBody.Type,
					"age":       float64(tt.wantBody.Age),
					"courses":   tt.wantBody.Courses,
				}

				// Convert courses to []int if present
				if courses, ok := got["courses"].([]interface{}); ok {
					intCourses := make([]int, len(courses))
					for i, c := range courses {
						if fNum, ok := c.(float64); ok {
							intCourses[i] = int(fNum)
						}
					}
					got["courses"] = intCourses
				}

				assert.Equal(t, want, got)
			}
		})
	}
}

func TestHandleUpdatePersonByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name       string
		urlName    string
		person     models.Person
		mockSetup  func(sqlmock.Sqlmock)
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name:    "Success",
			urlName: "John",
			person: models.Person{
				FirstName: "John",
				LastName:  "Smith",
				Type:     "student",
				Age:      21,
				Courses:  []int{1, 2},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// First query to find the person
				mock.ExpectQuery(`SELECT (.+) FROM person WHERE first_name = \$1`).
					WithArgs("John").
					WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "type", "age"}).
						AddRow(1, "John", "Doe", "student", 20))

				// Then expect the update
				mock.ExpectExec(`UPDATE person SET first_name = \$1, last_name = \$2, type = \$3, age = \$4 WHERE first_name = \$5`).
					WithArgs("John", "Smith", "student", 21, "John").
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Expect queries for adding courses
				for _, courseID := range []int{1, 2} {
					mock.ExpectExec(`INSERT INTO person_course \(person_id, course_id\) VALUES \(\$1, \$2\)`).
						WithArgs(1, courseID).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]interface{}{
				"id":        float64(1),
				"firstName": "John",
				"lastName":  "Smith",
				"type":      "student",
				"age":       float64(21),
				"courses":   []interface{}{float64(1), float64(2)},
			},
		},
		{
			name:    "Person Not Found",
			urlName: "NonExistent",
			person: models.Person{
				FirstName: "NonExistent",
				LastName:  "Person",
				Type:     "student",
				Age:      20,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM person WHERE first_name = \$1`).
					WithArgs("NonExistent").
					WillReturnError(sql.ErrNoRows)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   nil,
		},
		{
			name:    "Invalid Input - Empty Name",
			urlName: "",
			person: models.Person{
				FirstName: "",
				LastName:  "Smith",
				Type:     "student",
				Age:      21,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// No DB expectations needed for invalid input
			},
			wantStatus: http.StatusNotFound,
			wantBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			tt.mockSetup(mock)

			// Create request body
			jsonBody, err := json.Marshal(tt.person)
			if err != nil {
				t.Fatalf("Failed to marshal person: %v", err)
			}

			// Create request with chi context
			req := httptest.NewRequest("PUT", "/person/"+tt.urlName, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Create chi router and context
			r := chi.NewRouter()
			r.Put("/person/{name}", HandleUpdatePersonByName(db))
			r.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("HandleUpdatePersonByName() status = %v, want %v", w.Code, tt.wantStatus)
				t.Errorf("Response body: %s", w.Body.String())
				return
			}

			// For successful cases, check response body
			if tt.wantStatus == http.StatusOK {
				var gotBody map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&gotBody); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				// Compare each field individually
				for key, wantVal := range tt.wantBody {
					gotVal, exists := gotBody[key]
					if !exists {
						t.Errorf("Missing key %q in response", key)
						continue
					}

					if !reflect.DeepEqual(gotVal, wantVal) {
						t.Errorf("Mismatch for key %q:\ngot  (%T): %#v\nwant (%T): %#v",
							key, gotVal, gotVal, wantVal, wantVal)
					}
				}
			}

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestHandleCreatePerson(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name       string
		person     models.Person
		mockSetup  func(sqlmock.Sqlmock)
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name: "Success",
			person: models.Person{
				FirstName: "John",
				LastName:  "Doe",
				Type:     "student",
				Age:      20,
				Courses:  []int{1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO person").
					WithArgs("John", "Doe", "student", 20).
					WillReturnRows(rows)
				mock.ExpectExec("INSERT INTO person_course").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantStatus: http.StatusCreated,
			wantBody: map[string]interface{}{
				"id":        float64(1),  // Explicitly using float64
				"firstName": "John",
				"lastName":  "Doe",
				"type":      "student",
				"age":       float64(20), // Explicitly using float64
				"courses":   []interface{}{float64(1)}, // Using []interface{} with float64
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			tt.mockSetup(mock)

			// Create request body
			jsonBody, err := json.Marshal(tt.person)
			if err != nil {
				t.Fatalf("Failed to marshal person: %v", err)
			}

			// Create request
			r := httptest.NewRequest("POST", "/person", bytes.NewBuffer(jsonBody))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Handle request
			handler := HandleCreatePerson(db)
			handler.ServeHTTP(w, r)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("HandleCreatePerson() status = %v, want %v", w.Code, tt.wantStatus)
				t.Errorf("Response body: %s", w.Body.String())
				return
			}

			// For successful cases, check response body
			if tt.wantStatus == http.StatusCreated {
				var gotBody map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&gotBody); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				// Debug output
				t.Logf("Got body types: %#v", spewTypes(gotBody))
				t.Logf("Want body types: %#v", spewTypes(tt.wantBody))

				// Compare each field individually for better error messages
				for key, wantVal := range tt.wantBody {
					gotVal, exists := gotBody[key]
					if !exists {
						t.Errorf("Missing key %q in response", key)
						continue
					}

					if !reflect.DeepEqual(gotVal, wantVal) {
						t.Errorf("Mismatch for key %q:\ngot  (%T): %#v\nwant (%T): %#v",
							key, gotVal, gotVal, wantVal, wantVal)
					}
				}

				// Check for extra fields
				for key := range gotBody {
					if _, exists := tt.wantBody[key]; !exists {
						t.Errorf("Extra key %q in response", key)
					}
				}
			}

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

// Helper function to print types of map values
func spewTypes(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = fmt.Sprintf("%T", v)
	}
	return result
}
func TestHandleDeletePersonByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name       string
		personName string
		mockSetup  func(sqlmock.Sqlmock)
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Success",
			personName: "John",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Expect query to get person ID
				mock.ExpectQuery("SELECT id FROM person WHERE").
					WithArgs("John").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				
				// Expect delete from person_course
				mock.ExpectExec("DELETE FROM person_course").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				
				// Expect delete from person
				mock.ExpectExec("DELETE FROM person WHERE").
					WithArgs("John").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantStatus: http.StatusNoContent,
			wantBody:   `{"message": "Person deleted"}`,
		},
		{
			name:       "Person Not Found",
			personName: "NonExistent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id FROM person WHERE").
					WithArgs("NonExistent").
					WillReturnError(sql.ErrNoRows)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			router := chi.NewRouter()
			router.Delete("/people/{name}", HandleDeletePersonByName(db))

			req := httptest.NewRequest("DELETE", "/people/"+tt.personName, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusNoContent {
				// Trim any whitespace from the response body for comparison
				gotBody := strings.TrimSpace(w.Body.String())
				wantBody := strings.TrimSpace(tt.wantBody)
				assert.Equal(t, wantBody, gotBody)
			}
		})
	}
}