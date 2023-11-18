package isunippets

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

const SampleTableSchema = `
DROP TABLE IF EXISTS sample;
CREATE TABLE sample (
	string_value    VARCHAR(255) DEFAULT '',
	int_value 	    INTEGER      DEFAULT 0,
	timestamp_value TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
	bool_value 	    BOOLEAN      DEFAULT FALSE);
`

type SampleTableRow struct {
	String    string    `db:"string_value"`
	Int       int       `db:"int_value"`
	Timestamp time.Time `db:"timestamp_value"`
	Bool      bool      `db:"bool_value"`
}

func BulkInsert(db *sqlx.DB) error {
	var rows []SampleTableRow

	rows = append(rows, SampleTableRow{
		String:    "value1",
		Int:       1,
		Timestamp: time.Now(),
		Bool:      true,
	})

	rows = append(rows, SampleTableRow{
		String:    "value2",
		Int:       2,
		Timestamp: time.Unix(1700287644, 0),
		Bool:      false,
	})

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExec("INSERT INTO `sample` (`string_value`, `int_value`, `timestamp_value`, `bool_value`) VALUES (:string_value, :int_value, :timestamp_value, :bool_value)", rows)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func SelectNamedWithIn(db *sqlx.DB) ([]SampleTableRow, error) {
	query := "SELECT * FROM `sample` WHERE string_value IN (:stringValues) AND int_value = :intValue"

	input := map[string]interface{}{
		"stringValues": []string{"value1", "value2"},
		"intValue":     1,
	}

	query, args, err := sqlx.Named(query, input)
	if err != nil {
		return nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)

	var rows []SampleTableRow
	err = db.Select(&rows, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func SetUpSampleDatabase() (*sqlx.DB, error) {
	file, err := os.CreateTemp("", uuid.New().String())
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("sqlite3", file.Name())
	if err != nil {
		return nil, err
	}

	db.MustExec(SampleTableSchema)

	return db, nil
}
