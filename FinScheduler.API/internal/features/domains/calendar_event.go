package domains

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const calendarEventDateFormat = "2006-01-02"
const calendarEventTimeFormat = "15:04:05"
const calendarEventTimeFormatWithMicros = "15:04:05.999999"

const microsecondsPerSecond = int64(time.Second / time.Microsecond)
const microsecondsPerMinute = int64(time.Minute / time.Microsecond)
const microsecondsPerHour = int64(time.Hour / time.Microsecond)

type TimeOnly pgtype.Time

type CalendarEventColor string

const (
	Red    CalendarEventColor = "red"
	Blue   CalendarEventColor = "blue"
	Yellow CalendarEventColor = "yellow"
	Green  CalendarEventColor = "green"
	White  CalendarEventColor = "white"
	Orange CalendarEventColor = "orange"
	Violet CalendarEventColor = "violet"
)

func (timeOnly *TimeOnly) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*timeOnly = TimeOnly{}
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	parsedTime, err := parseCalendarEventTime(value, "time")
	if err != nil {
		return err
	}

	*timeOnly = parsedTime
	return nil
}

func (timeOnly TimeOnly) MarshalJSON() ([]byte, error) {
	if !timeOnly.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(formatCalendarEventTime(timeOnly))
}

func (timeOnly *TimeOnly) Scan(src any) error {
	var parsedTime pgtype.Time
	if err := parsedTime.Scan(src); err != nil {
		return err
	}

	*timeOnly = TimeOnly(parsedTime)
	return nil
}

func (timeOnly TimeOnly) Value() (driver.Value, error) {
	return pgtype.Time(timeOnly).Value()
}

func (calendarEventColor CalendarEventColor) IsValid() bool {
	switch calendarEventColor {
	case Red,
		Blue,
		Yellow,
		Green,
		White,
		Orange,
		Violet:
		return true
	default:
		return false
	}
}

type CalendarEvent struct {
	Id          uuid.UUID          `json:"id" db:"id"`
	Name        string             `json:"name" db:"name"`
	Description string             `json:"description" db:"description"`
	Color       CalendarEventColor `json:"color" db:"color"`
	Date        pgtype.Date        `json:"date" db:"date"`
}

type CalendarEventDateRangeFilter struct {
	From pgtype.Date
	To   pgtype.Date
}

type CalendarEventTriggerCreate struct {
	Time       TimeOnly `json:"time"`
	Commentary string   `json:"commentary"`
}

type CalendarEventTriggerUpdate struct {
	Id         *uuid.UUID `json:"id"`
	Time       TimeOnly   `json:"time"`
	Commentary string     `json:"commentary"`
}

type CalendarEventCreate struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Color       CalendarEventColor           `json:"color"`
	Date        pgtype.Date                  `json:"date"`
	Triggers    []CalendarEventTriggerCreate `json:"triggers"`
}

type CalendarEventUpdate struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Color       CalendarEventColor           `json:"color"`
	Date        pgtype.Date                  `json:"date"`
	Triggers    []CalendarEventTriggerUpdate `json:"triggers"`
}

func NewCalendarEventDateRangeFilter(r *http.Request) (CalendarEventDateRangeFilter, error) {
	queryParams := r.URL.Query()

	var from pgtype.Date
	fromRaw := queryParams.Get("from")
	if fromRaw != "" {
		parsedFrom, err := parseCalendarEventDate(fromRaw, "from")
		if err != nil {
			return CalendarEventDateRangeFilter{}, err
		}

		from = parsedFrom
	}

	var to pgtype.Date
	toRaw := queryParams.Get("to")
	if toRaw != "" {
		parsedTo, err := parseCalendarEventDate(toRaw, "to")
		if err != nil {
			return CalendarEventDateRangeFilter{}, err
		}

		to = parsedTo
	}

	return CalendarEventDateRangeFilter{
		From: from,
		To:   to,
	}, nil
}

func (calendarEvent *CalendarEventCreate) Validate() error {
	if !calendarEvent.Color.IsValid() {
		return fmt.Errorf("color is invalid")
	}

	if !calendarEvent.Date.Valid {
		return fmt.Errorf("date should be valid")
	}

	for _, trigger := range calendarEvent.Triggers {
		if !trigger.Time.Valid {
			return fmt.Errorf("trigger time should be valid")
		}
	}

	return nil
}

func (calendarEvent *CalendarEventUpdate) Validate() error {
	if !calendarEvent.Color.IsValid() {
		return fmt.Errorf("color is invalid")
	}

	if !calendarEvent.Date.Valid {
		return fmt.Errorf("date should be valid")
	}

	seen := make(map[uuid.UUID]struct{}, len(calendarEvent.Triggers))
	for _, trigger := range calendarEvent.Triggers {
		if !trigger.Time.Valid {
			return fmt.Errorf("trigger time should be valid")
		}

		if trigger.Id == nil {
			continue
		}

		if *trigger.Id == uuid.Nil {
			return fmt.Errorf("trigger id should not be nil")
		}

		if _, exists := seen[*trigger.Id]; exists {
			return fmt.Errorf("trigger id is duplicated: %s", trigger.Id.String())
		}

		seen[*trigger.Id] = struct{}{}
	}

	return nil
}

func (filter *CalendarEventDateRangeFilter) Validate() error {
	if !filter.From.Valid {
		return fmt.Errorf("from should be valid")
	}

	if !filter.To.Valid {
		return fmt.Errorf("to should be valid")
	}

	if filter.From.Time.After(filter.To.Time) {
		return fmt.Errorf("from should not be later than to")
	}

	return nil
}

func parseCalendarEventDate(value string, fieldName string) (pgtype.Date, error) {
	parsedDate, err := time.Parse(calendarEventDateFormat, value)
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("invalid %s value %q: %w", fieldName, value, err)
	}

	return pgtype.Date{
		Time:  parsedDate,
		Valid: true,
	}, nil
}

func parseCalendarEventTime(value string, fieldName string) (TimeOnly, error) {
	layouts := []string{calendarEventTimeFormat, calendarEventTimeFormatWithMicros}

	var lastErr error
	for _, layout := range layouts {
		parsedTime, err := time.Parse(layout, value)
		if err == nil {
			total := time.Duration(parsedTime.Hour())*time.Hour +
				time.Duration(parsedTime.Minute())*time.Minute +
				time.Duration(parsedTime.Second())*time.Second +
				time.Duration(parsedTime.Nanosecond()/int(time.Microsecond))*time.Microsecond

			return TimeOnly{
				Microseconds: int64(total / time.Microsecond),
				Valid:        true,
			}, nil
		}

		lastErr = err
	}

	return TimeOnly{}, fmt.Errorf("invalid %s value %q: %w", fieldName, value, lastErr)
}

func formatCalendarEventTime(value TimeOnly) string {
	remaining := value.Microseconds
	hours := remaining / microsecondsPerHour
	remaining = remaining % microsecondsPerHour
	minutes := remaining / microsecondsPerMinute
	remaining = remaining % microsecondsPerMinute
	seconds := remaining / microsecondsPerSecond
	microseconds := remaining % microsecondsPerSecond

	if microseconds == 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return fmt.Sprintf("%02d:%02d:%02d.%06d", hours, minutes, seconds, microseconds)
}
