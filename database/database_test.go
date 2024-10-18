package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jacob-tech-challenge/config"
)

func TestConnect_Success(t *testing.T) {
	// Step 1: Create a mock database
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error initializing mock DB: %v", err)
	}
	defer db.Close()

	// Step 2: Expect a successful ping
	mock.ExpectPing().WillReturnError(nil)

	// Step 3: Set up config
	cfg := config.Config{
		DB_Host:     "localhost",
		DB_Port:     5432,
		DB_User:     "user",
		DB_Password: "password",
		DB_Name:     "test_db",
	}

	// Step 4: Call the function with a mock OpenDBFunc
	conn, err := Connect(cfg, func(driverName, dataSourceName string) (*sql.DB, error) {
		return db, nil
	})

	// Step 5: Validate the connection
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if conn == nil {
		t.Fatalf("Expected valid connection, got nil")
	}

	// Step 6: Ensure expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unmet expectations: %v", err)
	}

	t.Logf("TestConnect_Success passed.")
}

func TestConnect_DBOpenFailure(t *testing.T) {
	// Set up config
	cfg := config.Config{
		DB_Host:     "localhost",
		DB_Port:     5432,
		DB_User:     "user",
		DB_Password: "password",
		DB_Name:     "test_db",
	}

	// Call the function with a mock OpenDBFunc that returns an error
	conn, err := Connect(cfg, func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, errors.New("DB connection failed")
	})

	// Validate
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if conn != nil {
		t.Errorf("Expected nil connection, got %v", conn)
	}
}

func TestConnect_NilDB(t *testing.T) {
	// Set up config
	cfg := config.Config{
		DB_Host:     "localhost",
		DB_Port:     5432,
		DB_User:     "user",
		DB_Password: "password",
		DB_Name:     "test_db",
	}

	// Call the function with a mock OpenDBFunc that returns nil
	conn, err := Connect(cfg, func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, nil
	})

	// Validate
	if err == nil || err.Error() != "db is nil" {
		t.Errorf("Expected 'db is nil' error, got %v", err)
	}
	if conn != nil {
		t.Errorf("Expected nil connection, got %v", conn)
	}
}

func TestConnect_PingFailure(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error initializing mock DB: %v", err)
	}
	defer db.Close()

	// Expect a ping failure
	mock.ExpectPing().WillReturnError(errors.New("ping failed"))

	// Set up config
	cfg := config.Config{
		DB_Host:     "localhost",
		DB_Port:     5432,
		DB_User:     "user",
		DB_Password: "password",
		DB_Name:     "test_db",
	}

	// Call the function with a mock OpenDBFunc
	conn, err := Connect(cfg, func(driverName, dataSourceName string) (*sql.DB, error) {
		return db, nil
	})

	// Validate
	if err == nil || err.Error() != "ping failed" {
		t.Errorf("Expected 'ping failed' error, got %v", err)
	}
	if conn != nil {
		t.Errorf("Expected nil connection, got %v", conn)
	}

	// Ensure expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %v", err)
	}
}
