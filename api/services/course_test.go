package services

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/jacob-tech-challenge/api/models"
)

func TestGetAllCourses(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Test case 1: Successful retrieval of multiple courses
	courses := []models.Course{
		{ID: 1, Name: "Course 1"},
		{ID: 2, Name: "Course 2"},
	}

	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, course := range courses {
		rows = rows.AddRow(course.ID, course.Name)
	}

	mock.ExpectQuery(`SELECT \* FROM "course"`).WillReturnRows(rows)

	retrievedCourses, err := GetAllCourses(db)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(retrievedCourses, courses) {
		t.Errorf("Returned courses do not match expected courses. Expected: %+v, Got: %+v", courses, retrievedCourses)
	}


	// Test case 2: No courses found
	mock.ExpectQuery(`SELECT \* FROM "course"`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	retrievedCourses, err = GetAllCourses(db)
	if err != nil {
		t.Errorf("Unexpected error when no courses are found: %v", err)
	}

	// CORRECTED ASSERTION:
	if len(retrievedCourses) != 0 {  // Check the length first
		t.Errorf("Expected empty slice, but got: %v", retrievedCourses)
	
	}


	// Test case 3: Database error
	mock.ExpectQuery(`SELECT \* FROM "course"`).WillReturnError(sql.ErrConnDone)

	_, err = GetAllCourses(db)
	if err == nil {
		t.Error("Expected an error, but got none")
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCourseByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testID := 1
	testCourse := models.Course{ID: testID, Name: "Test Course"}

	// Test successful retrieval
	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(testCourse.ID, testCourse.Name)
	mock.ExpectQuery(`SELECT \* FROM "course" WHERE id = \$1`).
		WithArgs(testID).
		WillReturnRows(rows)

	retrievedCourse, err := GetCourseByID(db, testID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(retrievedCourse, testCourse) {
		t.Errorf("Returned course does not match expected course. Expected: %+v, Got: %+v", testCourse, retrievedCourse)
	}

	// Test when no course is found
	mock.ExpectQuery(`SELECT \* FROM "course" WHERE id = \$1`).
		WithArgs(2).
		WillReturnError(sql.ErrNoRows)

	_, err = GetCourseByID(db, 2)
	if err == nil {
		t.Error("Expected sql.ErrNoRows, but got nil")
	}

	// Test with a database error
	mock.ExpectQuery(`SELECT \* FROM "course" WHERE id = \$1`).
		WithArgs(3).
		WillReturnError(sql.ErrConnDone)

	_, err = GetCourseByID(db, 3)
	if err == nil {
		t.Error("Expected an error, but got none")
	}


	// Ensure that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateCourse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testCourse := models.Course{
		ID:   1, // You might need to adjust this based on your model
		Name: "Updated Course Name",
	}

	// Define the expected query
	mock.ExpectExec(`UPDATE "course" SET name = \$1 WHERE id = \$2`).
		WithArgs(testCourse.Name, testCourse.ID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected, last insert ID 1 (doesn't matter for UPDATE)


	updatedCourse, err := UpdateCourse(db, testCourse.ID, testCourse)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert that the returned course is the same as the input course (after update)
	if !reflect.DeepEqual(updatedCourse, testCourse) {
		t.Errorf("Returned course does not match expected course. Expected: %+v, Got: %+v", testCourse, updatedCourse)

	}

	// Ensure that the query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//Test with an error
	mock.ExpectExec(`UPDATE "course" SET name = \$1 WHERE id = \$2`).
		WithArgs("Another Course Name", 2).
		WillReturnError(sql.ErrConnDone)

	_, err = UpdateCourse(db, 2, models.Course{ID:2, Name: "Another Course Name"})
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestCreateCourse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testCourse := models.Course{
		Name: "Test Course",
	}

	//Expect a new ID to be generated - we'll use 1 in the test, but this could be any valid ID.
	expectedID := int64(1)

	// Define the expected query
	mock.ExpectQuery(`INSERT INTO "course" \(\w+\) VALUES \(\$1\) RETURNING id`).
		WithArgs(testCourse.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	createdCourse, err := CreateCourse(db, testCourse)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert that the ID was set correctly and other fields are as expected
	if createdCourse.ID != int(expectedID) { //Casting from int64 to int, adjust if needed.
		t.Errorf("Returned course ID does not match expected ID. Expected: %d, Got: %d", expectedID, createdCourse.ID)
	}
	if createdCourse.Name != testCourse.Name {
		t.Errorf("Returned course name does not match expected name. Expected: %s, Got: %s", testCourse.Name, createdCourse.Name)
	}


	// Ensure that the query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}


	// Test with an error
	mock.ExpectQuery(`INSERT INTO "course" \(\w+\) VALUES \(\$1\) RETURNING id`).
		WithArgs("Error Course").
		WillReturnError(sql.ErrConnDone)

	_, err = CreateCourse(db, models.Course{Name: "Error Course"})
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteCourse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testID := 1

	//Expect a successful delete operation
	mock.ExpectExec(`DELETE FROM "course" WHERE id = \$1`).
		WithArgs(testID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

	err = DeleteCourse(db, testID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Ensure that the query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Test with an error - simulating no rows affected (though strictly not an error in SQL)
    mock.ExpectExec(`DELETE FROM "course" WHERE id = \$1`).
		WithArgs(2).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = DeleteCourse(db, 2)
	if err != nil {
		t.Errorf("Unexpected error when 0 rows are deleted: %v", err)
	}
	
    if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Test with a database error
	mock.ExpectExec(`DELETE FROM "course" WHERE id = \$1`).
		WithArgs(3).
		WillReturnError(sql.ErrConnDone)

	err = DeleteCourse(db, 3)
	if err == nil {
		t.Error("Expected an error, but got none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}