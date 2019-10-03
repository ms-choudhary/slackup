package util

import (
	"database/sql"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func SetupDatabaseFrom(dbFile string, sqlFiles ...string) error {
	os.Remove(dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}

	for _, f := range sqlFiles {
		sql, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(sql))
		if err != nil {
			return err
		}
	}
	return nil
}
