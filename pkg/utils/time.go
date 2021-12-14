package utils

import "time"

type timeLayout string

const (
	defaultTimeLayout timeLayout = "2006-01-02 15:04:05"
	timeLayoutTZ      timeLayout = "2006-01-02T15:04:05Z"
)

func CurrentTime() string {
	return current(defaultTimeLayout)
}

func CurrentTZTime() string {
	return current(timeLayoutTZ)
}

func current(layout timeLayout) string {
	return time.Now().Format(string(layout))
}

func ParseTime(value string) time.Time {
	return parse(defaultTimeLayout, value)
}

func ParseTimeInLocation(value string, loc *time.Location) time.Time {
	t, _ := parseInLocation(defaultTimeLayout, value, loc)
	return t
}

func ParseTimeTz(value string) time.Time {
	return parse(timeLayoutTZ, value)
}

func ParseTimeTzInLocation(value string, loc *time.Location) time.Time {
	t, _ := parseInLocation(timeLayoutTZ, value, loc)
	return t
}

func parseInLocation(layout timeLayout, value string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(string(layout), value, loc)
}

func parse(layout timeLayout, value string) time.Time {
	t, _ := time.Parse(string(layout), value)
	return t
}

func FormatTime(t time.Time) string {
	return format(defaultTimeLayout, t)
}

func format(layout timeLayout, t time.Time) string {
	return t.Format(string(layout))
}
