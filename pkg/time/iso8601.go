package timepkg

import "time"

func TimeToISO8601(v time.Time) string {
	return v.Format("2006-01-02T15:04:05Z")
}
