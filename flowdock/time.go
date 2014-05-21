package flowdock

import (
	"strconv"
	"time"
)

// Time represents a Flowdock time stamp which is milliseconds since Epoch
// based
type Time struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface. The time is
// expected to be an intger representing milliseconds since Epoch.
func (t *Time) UnmarshalJSON(b []byte) error {
	result, err := strconv.ParseInt(string(b), 0, 64)

	if err != nil {
		return err
	}

	// convert the unix epoch to a Time object
	t.Time = time.Time(time.Unix(result/1000, 0))

	return nil
}
