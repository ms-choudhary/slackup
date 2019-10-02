package store

import (
	"database/sql"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupDatabaseFrom(dbFile string, sqlFiles ...string) error {
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

func TestGetChannel(t *testing.T) {
	cases := []struct {
		fixture          string
		project, channel string
		returnId         int
		name             string
	}{
		{
			fixture:  "../../testdata/db/empty.sql",
			project:  "scribpox",
			channel:  "ops-incident",
			returnId: 1,
			name:     "ChannelDoesNotExists",
		},
		{
			fixture:  "../../testdata/db/channel.sql",
			project:  "scripbox",
			channel:  "ops-incident",
			returnId: 1,
			name:     "ChannelExists",
		},
		{
			fixture:  "../../testdata/db/channel.sql",
			project:  "scribpox",
			channel:  "general",
			returnId: 2,
			name:     "NewChannel",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			err := setupDatabaseFrom("test.db", "../../scripts/migration.sql", tc.fixture)
			if err != nil {
				t.Fatalf("did not expect error: %v", err)
			}

			db, err := Init("test.db")
			if err != nil {
				t.Fatalf("did not expect error: %v", err)
			}

			returnId, err := GetChannel(db, tc.project, tc.channel)
			if err != nil {
				t.Fatalf("did not expect error: %v", err)
			}

			if returnId != tc.returnId {
				t.Fatalf("expected %v got %v", tc.returnId, returnId)
			}
		})

	}
}
