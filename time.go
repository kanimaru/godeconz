package godeconz

import "time"

// deCONZ timestamps come in several shapes: "lastupdated" has milliseconds and no zone
// ("2024-05-01T10:00:01.123"), "lastseen" has minute resolution with Z ("2024-05-01T10:00Z").
// All of them are UTC.
var timeLayouts = []string{
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04Z",
	time.RFC3339Nano,
}

// ParseTime parses a deCONZ timestamp as UTC. Empty strings and "none" return ok == false.
func ParseTime(s string) (t time.Time, ok bool) {
	if s == "" || s == "none" {
		return time.Time{}, false
	}
	for _, layout := range timeLayouts {
		if parsed, err := time.ParseInLocation(layout, s, time.UTC); err == nil {
			return parsed.UTC(), true
		}
	}
	return time.Time{}, false
}
