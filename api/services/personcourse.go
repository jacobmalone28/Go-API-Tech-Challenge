package services

import (
	"context"
	"database/sql"
	"fmt"
)

// AddPersonToCourse adds a person to multiple courses, handling potential errors for individual courses.
func AddPersonToCourse(db *sql.DB, personID int, courseIDs []int) error {
    ctx := context.Background()
    var errors []error

    for _, id := range courseIDs {
		// make sure the course exists
		var courseID int
		err := db.QueryRowContext(ctx, `SELECT id FROM course WHERE id = $1`, id).Scan(&courseID)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to add course %d: course not found", id))
			continue
		}

		// make sure the person exists
		var personID int
		err = db.QueryRowContext(ctx, `SELECT id FROM person WHERE id = $1`, personID).Scan(&personID)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to add course %d: person not found", id))
			continue
		}



        _, err = db.ExecContext(
            ctx,
            `INSERT INTO person_course (person_id, course_id) VALUES ($1, $2)`,
            personID, id,
        )
        if err != nil {
            // Collect errors instead of just continuing
            errors = append(errors, fmt.Errorf("failed to add course %d: %w", id, err))
        }
    }

    // If there were any errors, return them combined
    if len(errors) > 0 {
        return fmt.Errorf("failed to add some courses: %v", errors)
    }

    return nil
}