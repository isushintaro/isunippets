package isunippets

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestBulkInsert(t *testing.T) {
	assert := assert.New(t)

	file, err := os.CreateTemp("", "bulk_insert_test")
	assert.NoError(err)

	db, err := sqlx.Connect("sqlite3", file.Name())
	assert.NoError(err)

	var schema = `
DROP TABLE IF EXISTS bulk;
CREATE TABLE bulk (
	string_value    VARCHAR(255) DEFAULT '',
	int_value 	    INTEGER      DEFAULT 0,
	timestamp_value TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
	bool_value 	    BOOLEAN      DEFAULT FALSE);
`
	db.MustExec(schema)

	ts := time.Now().Unix()
	assert.NoError(BulkInsert(db))

	row := BulkTableRow{}
	err = db.Get(&row, "SELECT * FROM bulk where string_value = 'value1'")
	assert.NoError(err)

	assert.Equal("value1", row.String)
	assert.Equal(1, row.Int)
	assert.Equal(true, row.Bool)
	assert.GreaterOrEqual(row.Timestamp.Unix(), ts)
}
