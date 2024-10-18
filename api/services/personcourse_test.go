package services

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAddPersonToCourse(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	
	tests := map[string]struct{
		personID int
		courseID int
		expectedErr string
	}{
		"success": {
			personID: 1,
			courseID: 1,
		},
		"person not found": {
			personID: 0,
			courseID: 1,
			expectedErr: "failed to add some courses: [failed to add course 1: course not found]",
		},
		"course not found": {
			personID: 1,
			courseID: 0,
			expectedErr: "failed to add some courses: [failed to add course 0: course not found]",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T){
			mock.ExpectQuery(`SELECT id FROM course WHERE id = \$1`).WithArgs(tc.courseID).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			mock.ExpectQuery(`SELECT id FROM person WHERE id = \$1`).WithArgs(tc.personID).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			mock.ExpectExec(`INSERT INTO person_course \(person_id, course_id\) VALUES \(\$1, \$2\)`).WithArgs(tc.personID, tc.courseID).WillReturnResult(sqlmock.NewResult(1, 1))

			err := AddPersonToCourse(db, tc.personID, []int{tc.courseID})

			if err != nil && err.Error() != tc.expectedErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}

		})
	}
}