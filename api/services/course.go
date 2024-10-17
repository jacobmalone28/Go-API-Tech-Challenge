package services

import (
	"context"
	"database/sql"
	"log"

	"github.com/jacob-tech-challenge/api/models"
)

// GetAllCourses returns all courses
func GetAllCourses(db *sql.DB) ([]models.Course, error) {
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, `SELECT * FROM "course"`)
	if err != nil {
		return []models.Course{}, err // Return early if there's an error in QueryContext
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// Log the close error, but don't mask the original error.
			log.Printf("Error closing rows: %v", closeErr)
		}
	}()


	var courses []models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	if err := rows.Err(); err != nil {
		return nil, err // Check for errors after iteration
	}
	return courses, nil
}

// GetCourseByID returns a course by id
func GetCourseByID(db *sql.DB, id int) (models.Course, error) {

	ctx := context.Background()

	var course models.Course

	if err := db.QueryRowContext(
		ctx,
		`SELECT * FROM "course" WHERE id = $1`,
		id,
	).Scan(&course.ID, &course.Name); err != nil {
		return models.Course{}, err
	}
	return course, nil
}

// UpdateCourse updates a course
func UpdateCourse(db *sql.DB, id int, course models.Course) (models.Course, error) {
	
	ctx := context.Background()

	_, err := db.ExecContext(
		ctx,
		`UPDATE "course" SET name = $1 WHERE id = $2`,
		course.Name, id,
	)
	if err != nil {
		return models.Course{}, err
	}
	// return the updated course
	return course, nil
}

// CreateCourse creates a course
func CreateCourse(db *sql.DB, course models.Course) (models.Course, error) {
	
	ctx := context.Background()

	err := db.QueryRowContext(
		ctx,
		`INSERT INTO "course" (name) VALUES ($1) RETURNING id`,
		course.Name,
	).Scan(&course.ID)
	if err != nil {
		return models.Course{}, err
	}
	return course, nil
}

// DeleteCourse deletes a course
func DeleteCourse(db *sql.DB, id int) error {
	
	ctx := context.Background()

	_, err := db.ExecContext(
		ctx,
		`DELETE FROM "course" WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
