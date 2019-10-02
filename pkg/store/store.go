package store

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ms-choudhary/slackup/pkg/api"
)

type Filter struct {
	User string
	Text string
	Ts   string
}

func Init(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open sqlite3 db: %s %v\n", dbPath, err)
	}
	return db, nil
}

func createChannel(db *sql.DB, project, channel string) (int, error) {
	stmt, err := db.Prepare("INSERT INTO channel(project_name, channel_name) VALUES(?, ?)")

	res, err := stmt.Exec(project, channel)
	if err != nil {
		return -1, fmt.Errorf("cannot create channel: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("cannot get channel's last inserted id: %v", err)
	}
	return int(id), nil
}

func channelExists(db *sql.DB, project, channel string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM channel WHERE project_name = ? AND channel_name = ?)", project, channel).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check if channel exists: %v", err)
	}
	return exists, nil
}

// Create project, channel if does not exists, else return channel id
func GetChannel(db *sql.DB, project, channel string) (int, error) {
	exists, err := channelExists(db, project, channel)
	if err != nil {
		return -1, err
	}

	if !exists {
		id, err := createChannel(db, project, channel)
		if err != nil {
			return -1, err
		}
		return id, nil
	}

	var id int
	err = db.QueryRow("SELECT ID FROM channel WHERE project_name = ? AND channel_name = ?", project, channel).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("failed to get channel id: %v", err)
	}
	return id, nil
}

// Insert messages for a channel
func Insert(db *sql.DB, channel int, messages []api.Message) error {

	stmt, err := db.Prepare("INSERT INTO message(user, text, ts, channel_id, parent_id) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, msg := range messages {
		res, err := stmt.Exec(msg.User, msg.Text, msg.Ts, channel, -1)
		if err != nil {
			return err
		}

		msgId, err := res.LastInsertId()
		if err != nil {
			return err
		}

		for _, c := range msg.Comments {
			_, err := stmt.Exec(c.User, c.Text, c.Ts, channel, msgId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getComments(stmt *sql.Stmt, channel, parentId int) ([]api.Message, error) {
	rows, err := stmt.Query(channel, parentId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []api.Message{}
	for rows.Next() {
		var (
			id             int
			user, text, ts string
		)
		err := rows.Scan(&user, &text, &ts, &id)
		if err != nil {
			return nil, err
		}

		comments = append(comments, api.Message{User: user, Text: text, Ts: ts})
	}
	return comments, nil
}

// Query a channel by filter
// TODO filter is ignored right now
func Query(db *sql.DB, channel int, filter Filter) ([]api.Message, error) {
	stmt, err := db.Prepare("SELECT user, text, ts, id FROM message WHERE channel_id = ? AND parent_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	msgRows, err := stmt.Query(channel, -1)
	if err != nil {
		return nil, err
	}

	defer msgRows.Close()

	messages := []api.Message{}
	for msgRows.Next() {
		var (
			id             int
			user, text, ts string
		)
		err := msgRows.Scan(&user, &text, &ts, &id)
		if err != nil {
			return nil, err
		}

		comments, err := getComments(stmt, channel, id)
		if err != nil {
			return nil, err
		}

		messages = append(messages, api.Message{User: user, Text: text, Ts: ts, Comments: comments})
	}

	return messages, nil
}
