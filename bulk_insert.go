package isunippets

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type BulkTableRow struct {
	String    string    `db:"string_value"`
	Int       int       `db:"int_value"`
	Timestamp time.Time `db:"timestamp_value"`
	Bool      bool      `db:"bool_value"`
}

func BulkInsert(db *sqlx.DB) error {
	var params []BulkTableRow

	params = append(params, BulkTableRow{
		String:    "value1",
		Int:       1,
		Timestamp: time.Now(),
		Bool:      true,
	})

	params = append(params, BulkTableRow{
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

	_, err = tx.NamedExec("INSERT INTO `bulk` (`string_value`, `int_value`, `timestamp_value`, `bool_value`) VALUES (:string_value, :int_value, :timestamp_value, :bool_value)", params)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
