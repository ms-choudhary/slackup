package store

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ms-choudhary/slackup/pkg/api"
)

var (
	migrationFile   = "../../scripts/migration.sql"
	fixturesBaseDir = "../../testdata/db"
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

func expectNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("did not expect error: %v", err)
	}
}

func isEqualMessage(a, b api.Message) bool {
	if a.User != b.User || a.Text != b.Text || a.Ts != b.Ts || len(a.Comments) != len(b.Comments) {
		return false
	}

	for i, c := range a.Comments {
		if !isEqualMessage(c, b.Comments[i]) {
			return false
		}
	}

	return true
}

func isEqualMessages(a, b []api.Message) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if !isEqualMessage(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestInsertQuery(t *testing.T) {
	cases := []struct {
		name     string
		fixture  string
		channel  int
		messages []api.Message
	}{
		{
			name:    "WYSIWYG",
			fixture: filepath.Join(fixturesBaseDir, "channel.sql"),
			channel: 1,
			messages: []api.Message{
				api.Message{
					User: "mohit",
					Text: "hello, world",
					Comments: []api.Message{
						api.Message{
							User: "bot",
							Text: "howdy",
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := setupDatabaseFrom("test.db", migrationFile, tc.fixture)
			expectNoError(t, err)

			db, err := Init("test.db")
			expectNoError(t, err)

			err = Insert(db, tc.channel, tc.messages)
			expectNoError(t, err)

			retMsgs, err := Query(db, tc.channel, Filter{})

			if !isEqualMessages(tc.messages, retMsgs) {
				t.Fatalf("expected %v got %v", tc.messages, retMsgs)
			}
		})
	}
}

func TestGetChannel(t *testing.T) {
	cases := []struct {
		fixture          string
		project, channel string
		returnId         int
		name             string
	}{
		{
			fixture:  filepath.Join(fixturesBaseDir, "empty.sql"),
			project:  "scribpox",
			channel:  "ops-incident",
			returnId: 1,
			name:     "ChannelDoesNotExists",
		},
		{
			fixture:  filepath.Join(fixturesBaseDir, "channel.sql"),
			project:  "scripbox",
			channel:  "ops-incident",
			returnId: 1,
			name:     "ChannelExists",
		},
		{
			fixture:  filepath.Join(fixturesBaseDir, "channel.sql"),
			project:  "scribpox",
			channel:  "general",
			returnId: 2,
			name:     "NewChannel",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			err := setupDatabaseFrom("test.db", migrationFile, tc.fixture)
			expectNoError(t, err)

			db, err := Init("test.db")
			expectNoError(t, err)

			returnId, err := GetChannel(db, tc.project, tc.channel)
			expectNoError(t, err)

			if returnId != tc.returnId {
				t.Fatalf("expected %v got %v", tc.returnId, returnId)
			}
		})

	}
}
