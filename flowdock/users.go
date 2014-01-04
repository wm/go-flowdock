package flowdock

type User struct {
	id            *string
	nick          *string
	name          *string
	email         *string
	avatar        *string
	status        *string
	disabled      *bool
	last_activity *Time
	last_ping     *Time
}
