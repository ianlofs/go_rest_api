package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	Conn *sql.DB
}

func New() (DB, error) {
  dbURL := os.Getenv("DATABASE_URL")

  if dbURL == "" {
    dbURL = fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
      os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PSSWRD"), os.Getenv("DB_PORT"))
  }
	conn, err := sql.Open("postgres", dbURL)
	return DB{conn}, err
}

func (db DB) Close() error {
	return db.Conn.Close()
}

func (db DB) SetMaxOpenConns(n int) {
	db.Conn.SetMaxOpenConns(n)
}

func (db DB) Insert()  {
  /*
	log.Println("Beginning transaction...")
	// Begin transaction. Required for bulk insert
	txn, err := db.conn.Begin()
	if err != nil {
		return
	}

	// Prepare bulk insert statement
	stmt, err := txn.Prepare(pq.CopyIn("data_raw", "serial", "type", "data", "time", "device"))

	// Cleanup either when done or in the case of an error
	defer func() {
		log.Println("Closing off transaction...")

		if stmt != nil {
			// Flush buffer
			if _, eerr := stmt.Exec(); eerr != nil {
				if err == nil {
					err = eerr
				}
			}

			// Close prepared statement
			if cerr := stmt.Close(); cerr != nil {
				if err == nil {
					err = cerr
				}
			}
		}

		// Rollback transaction on error
		if err != nil {
			txn.Rollback()
			log.Println("Transaction rolled back")
			return
		}

		// Commit transaction
		err = txn.Commit()

		log.Println("Transaction closed")
	}()

	// Check for error from preparing statement
	if err != nil {
		return
	}

	for {
		var row *decoders.DataPoint
		row, err = iter()
		if row == nil || err != nil {
			break
		}

		if constants.Verbose {
			log.Println("Data:", row.Data)
			log.Println("Time:", row.Time)
		}

		// Insert data. This is buffered.
		if _, err = stmt.Exec(row.Serial, string(row.Type), row.Data, row.Time, row.Device); err != nil {
			break
		}
	}
	return*/
}
