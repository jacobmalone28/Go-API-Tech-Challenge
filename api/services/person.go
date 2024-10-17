package services

import (
	"context"
	"database/sql"

	"github.com/jacob-tech-challenge/api/models"
)

// GetAllPeople returns all people, if query parameters are provided, it filters the results
func GetAllPeople(db *sql.DB, name string, age int) ([]models.Person, error) {
	ctx := context.Background()

	var rows *sql.Rows
	var err error

	rows, err = db.QueryContext(ctx, `SELECT * FROM person WHERE ($1 = '' OR first_name = $1) AND ($2 = 0 OR age = $2)`, name, age)


	if err != nil {
		return []models.Person{}, err
	}
	defer rows.Close()


	var people []models.Person
	for rows.Next() {
		var person models.Person
		if err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &person.Type, &person.Age); err != nil {
			return nil, err
		}
		person.Courses, err = GetCoursesByPersonID(db, person.ID)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}
	return people, nil
}

// GetPersonByName returns a person by name
func GetPersonByName(db *sql.DB, name string) (models.Person, error) {
	ctx := context.Background()

	var person models.Person
	err := db.QueryRowContext(ctx, `SELECT * FROM person WHERE first_name = $1`, name).Scan(&person.ID, &person.FirstName, &person.LastName, &person.Type, &person.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Person{}, nil
		}
		return models.Person{}, err
	}
	person.Courses, err = GetCoursesByPersonID(db, person.ID)
	if err != nil {
		return models.Person{}, err
	}
	return person, nil
}

// UpdatePersonByName updates a person by name
func UpdatePersonByName(db *sql.DB, name string, person models.Person) (models.Person, error) {
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `UPDATE person SET first_name = $1, last_name = $2, type = $3, age = $4 WHERE first_name = $5`, person.FirstName, person.LastName, person.Type, person.Age, name)
	if err != nil {
		return models.Person{}, err
	}
	return GetPersonByName(db, person.FirstName)
}

// CreatePerson creates a person
func CreatePerson(db *sql.DB, person models.Person) (models.Person, error) {
    ctx := context.Background()

    // Start a transaction
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return models.Person{}, err
    }
    defer tx.Rollback()

    // Insert person
    err = tx.QueryRowContext(ctx, 
        `INSERT INTO person (first_name, last_name, type, age) VALUES ($1, $2, $3, $4) RETURNING id`,
        person.FirstName, person.LastName, person.Type, person.Age).Scan(&person.ID)
    if err != nil {
        return models.Person{}, err
    }

    // Insert course associations
    for _, courseID := range person.Courses {
        _, err = tx.ExecContext(ctx,
            `INSERT INTO person_course (person_id, course_id) VALUES ($1, $2)`,
            person.ID, courseID)
        if err != nil {
            return models.Person{}, err
        }
    }

    // Commit transaction
    if err = tx.Commit(); err != nil {
        return models.Person{}, err
    }

    return person, nil
}

// DeletePersonByName deletes a person by name
func DeletePersonByName(db *sql.DB, name string) error {
	ctx := context.Background()

	// grab the person id
	var person models.Person
	err := db.QueryRowContext(ctx, `SELECT id FROM person WHERE first_name = $1`, name).Scan(&person.ID)

	if err != nil {
		return err
	}

	// delete the person from the person_course table
	_, err = db.ExecContext(ctx, `DELETE FROM person_course WHERE person_id = $1`, person.ID)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, `DELETE FROM person WHERE first_name = $1`, name)
	if err != nil {
		return err
	}	
	return nil
}

// GetCoursesByPersonID returns all course ids for a person
func GetCoursesByPersonID(db *sql.DB, personID int) ([]int, error) {
	
	ctx := context.Background()

	rows, err := db.QueryContext(
		ctx,
		`SELECT c.id, c.name FROM "course" c JOIN "person_course" pc ON c.id = pc.course_id WHERE pc.person_id = $1`,
		personID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []int
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name); err != nil {
			return nil, err
		}
		courses = append(courses, course.ID)
	}
	return courses, nil
}

