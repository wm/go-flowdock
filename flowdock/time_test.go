package flowdock

import (
	"testing"
	"time"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	flowdockTime := Time{}
	json := []byte("1385546251160")
	err := flowdockTime.UnmarshalJSON(json)
	if err != nil {
		t.Errorf("Time.UnmarshalJSON returned error: %v", err)
	}

	want := time.Date(2013, time.November, 27, 9, 57, 31, 0, time.UTC)
	if flowdockTime.Local() != want.Local() {
		t.Errorf("Time.UnmarshalJSON set time to %v, wanted %v", flowdockTime.Local(), want.Local())
	}
}
