package services

import (
	"context"
	"database/sql"

	"github.com/jacob-tech-challenge/api/models"
)

// GetAllCourses returns all courses
func GetAllCourses(db *sql.DB) ([]models.Course, error) {

	ctx := context.Background()

	rows, err := db.QueryContext(
		ctx,
		`SELECT * FROM "course"`,
	)
	if err != nil {
		return []models.Course{}, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name); err != nil {
			return nil, err
		}
		courses = append(courses, course)
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
func UpdateCourse(db *sql.DB, course models.Course) (models.Course, error) {
	
	ctx := context.Background()

	_, err := db.ExecContext(
		ctx,
		`UPDATE "course" SET name = $1 WHERE id = $2`,
		course.Name, course.ID,
	)
	if err != nil {
		return models.Course{}, err
	}
	// return the updated course
	return course, nil
}