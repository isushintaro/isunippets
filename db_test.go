package isunippets

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBulkInsert(t *testing.T) {
	assert := assert.New(t)

	db, err := SetUpSampleDatabase()
	assert.NoError(err)
	defer db.Close()

	ts := time.Now().Unix()
	assert.NoError(BulkInsert(db))

	row := SampleTableRow{}
	err = db.Get(&row, "SELECT * FROM sample where string_value = 'value1'")
	assert.NoError(err)

	assert.Equal("value1", row.String)
	assert.Equal(1, row.Int)
	assert.Equal(true, row.Bool)
	assert.GreaterOrEqual(row.Timestamp.Unix(), ts)
}

func TestSelectNamedWithIn(t *testing.T) {
	assert := assert.New(t)

	db, err := SetUpSampleDatabase()
	assert.NoError(err)
	defer db.Close()

	_, err = db.NamedExec("INSERT INTO `sample` (`string_value`, `int_value`, `timestamp_value`, `bool_value`) VALUES (:string_value, :int_value, :timestamp_value, :bool_value)", []SampleTableRow{
		{
			String: "value1",
			Int:    1,
		},
		{
			String: "value2",
			Int:    1,
		},
		{
			String: "value3",
			Int:    1,
		},
		{
			String: "value1",
			Int:    2,
		},
	})
	assert.NoError(err)

	rows, err := SelectNamedWithIn(db)
	assert.NoError(err)
	assert.Len(rows, 2)
}
