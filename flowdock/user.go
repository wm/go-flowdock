package flowdock

import (
	"time"
)

type User struct {
	id             string
	nick           string
	name           string
	email          string
	avatar         string
	status         string
	disabled       bool
	last_activity  time.Time
	last_ping      time.Time
}
