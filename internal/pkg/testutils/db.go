package testutils

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

const (
	initSchemaQuery = `CREATE TABLE positions (keyword text, position integer, domain text, url text, volume integer, results integer, cpc float, updated datetime, primary key (domain, url, keyword));`

	dataInsertQuery = `INSERT INTO positions (keyword, position , domain, url, volume, results, cpc, updated) VALUE 
("test1", 1, "ulmart.ru", "http://ulmart.ru/test1", 43, 40000, 3.22, 1495248847),
("test2", 2, "ulmart.ru", "http://ulmart.ru/test2", 55, 40000, 1.22, 1495248847),
("test3", 3, "ulmart.ru", "http://ulmart.ru/test3", 76, 40000, 2.22, 1495248847),
("test4", 7, "non-ulmart.ru", "http://nonulmart.ru/tests/2", 32, 40000, 5.22, 1495248847),
("test4", 11, "non-ulmart.ru", "http://nonulmart.ru/tests", 65, 40000, 5.22, 1495248847);`

	dropTableQuery = `DROP TABLE positions`
)

func PrepareDB(t *testing.T, conn *sqlx.DB) {
	_, err := conn.Exec(initSchemaQuery)
	assert.NoError(t, err)
	_, err = conn.Exec(dataInsertQuery)
	assert.NoError(t, err)
}

func TeardownDB(t *testing.T, conn *sqlx.DB) {
	_, err := conn.Exec(dropTableQuery)
	assert.NoError(t, err)
}
