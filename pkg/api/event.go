package api

import (
	"time"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/middleware"
	m "github.com/grafana/grafana/pkg/models"
)

func GetEvents(c *middleware.Context, query m.GetEventsQuery) {
	query.AccountId = c.AccountId

	if query.End == 0 {
		query.End = time.Now().Unix() * 1000
	}
	if query.Start == 0 {
		query.Start = query.End - (60 * 60 * 1000) //1hour
	}

	if query.Size == 0 {
		query.Size = 10
	}

	if err := bus.Dispatch(&query); err != nil {
		c.JsonApiErr(500, "Failed to query events", err)
		return
	}

	c.JSON(200, query.Result)
}