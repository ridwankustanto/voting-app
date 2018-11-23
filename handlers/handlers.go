package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"voting-app/models"

	"github.com/labstack/echo"
)

// Handler for get all polls

// GetPolls to handle response get poll
func GetPolls(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, models.GetPolls(db))
	}
}

// UpdatePoll to handle update poll
func UpdatePoll(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var poll models.Poll

		c.Bind(&poll)

		index, _ := strconv.Atoi(c.Param("index"))

		id, err := models.UpdatePoll(db, index, poll.Name, poll.Upvotes, poll.Downvotes)

		if err == nil {
			return c.JSON(http.StatusCreated, id)
		}

		return err
	}
}
