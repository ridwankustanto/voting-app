package models

import (
	"database/sql"

	pusher "github.com/pusher/pusher-http-go"
)

var client = pusher.Client{
	AppId:   "655961",
	Key:     "50aa1a6ab900a700ae4f",
	Secret:  "f8035f7adb4c1a8ea655",
	Cluster: "ap1",
	Secure:  true,
}

// Poll struct for json object poll
type Poll struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Topic     string `json:"topic"`
	Src       string `json:"src"`
	Upvotes   int    `json:"upvotes"`
	Downvotes int    `json:"downvotes"`
}

// PollCollection struct for retrieving all poll
type PollCollection struct {
	Polls []Poll `json:"items"`
}

// GetPolls for query to get all poll
func GetPolls(db *sql.DB) PollCollection {
	sql := "SELECT * FROM polls"

	rows, err := db.Query(sql)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	result := PollCollection{}

	for rows.Next() {
		poll := Poll{}

		err2 := rows.Scan(&poll.ID, &poll.Name, &poll.Topic, &poll.Src, &poll.Upvotes, &poll.Downvotes)

		if err2 != nil {
			panic(err2)
		}

		result.Polls = append(result.Polls, poll)
	}

	return result
}

// UpdatePoll for query update upvotes or downvotes the poll
func UpdatePoll(db *sql.DB, index int, name string, upvotes int, downvotes int) (int64, error) {
	sql := "UPDATE polls SET (upvotes, downvotes) = (?, ?) WHERE id = ?"

	// Create a prepared SQL statement
	stmt, err := db.Prepare(sql)

	// Exit if we get an error
	if err != nil {
		panic(err)
	}

	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Replace the '?' in our prepared statement with 'upvotes, downvotes, index'
	result, err2 := stmt.Exec(upvotes, downvotes, index)

	// Exit if we get an error
	if err2 != nil {
		panic(err2)
	}

	pollUpdate := Poll{
		ID:        index,
		Name:      name,
		Upvotes:   upvotes,
		Downvotes: downvotes,
	}

	client.Trigger("poll-channel", "poll-update", pollUpdate)
	return result.RowsAffected()
}
