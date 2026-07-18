package domains

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCalendarEventDateRangeFilter_ShouldParseSupportedFields(t *testing.T) {
	// Arrange
	requestURL := "/calendar-events?from=2026-07-17&to=2026-07-18"
	request := httptest.NewRequest("GET", requestURL, nil)

	// Act
	filter, err := NewCalendarEventDateRangeFilter(request)

	// Assert
	require.NoError(t, err)
	assert.True(t, filter.From.Valid)
	assert.True(t, filter.To.Valid)
	assert.Equal(t, "2026-07-17", filter.From.Time.Format("2006-01-02"))
	assert.Equal(t, "2026-07-18", filter.To.Time.Format("2006-01-02"))
}

func TestNewCalendarEventDateRangeFilter_ShouldReturnZeroValueWhenQueryIsEmpty(t *testing.T) {
	// Arrange
	requestURL := "/calendar-events"
	request := httptest.NewRequest("GET", requestURL, nil)

	// Act
	filter, err := NewCalendarEventDateRangeFilter(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, CalendarEventDateRangeFilter{}, filter)
}

func TestNewCalendarEventDateRangeFilter_ShouldReturnErrorOnInvalidQueryParam(t *testing.T) {
	// Arrange
	requestURL := "/calendar-events?from=bad-date&to=2026-07-18"
	request := httptest.NewRequest("GET", requestURL, nil)

	// Act
	filter, err := NewCalendarEventDateRangeFilter(request)

	// Assert
	require.Error(t, err)
	assert.Equal(t, CalendarEventDateRangeFilter{}, filter)
	assert.Contains(t, err.Error(), `invalid from value "bad-date"`)
}

func TestCalendarEventCreateUnmarshalJSON_ShouldParseSupportedFields(t *testing.T) {
	// Arrange
	payload := []byte(`{"name":"Bills","description":"Rent","color":"red","date":"2026-07-17","triggers":[{"time":"08:30:00","commentary":"Morning reminder"},{"time":"20:00:00","commentary":"Evening reminder"}]}`)

	// Act
	var create CalendarEventCreate
	err := json.Unmarshal(payload, &create)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Bills", create.Name)
	assert.Equal(t, "Rent", create.Description)
	assert.Equal(t, "red", create.Color)
	assert.True(t, create.Date.Valid)
	assert.Equal(t, "2026-07-17", create.Date.Time.Format("2006-01-02"))
	require.Len(t, create.Triggers, 2)
	assert.True(t, create.Triggers[0].Time.Valid)
	assert.Equal(t, int64(30600000000), create.Triggers[0].Time.Microseconds)
	assert.Equal(t, "Morning reminder", create.Triggers[0].Commentary)
	assert.Equal(t, int64(72000000000), create.Triggers[1].Time.Microseconds)
	assert.Equal(t, "Evening reminder", create.Triggers[1].Commentary)
}

func TestCalendarEventCreateUnmarshalJSON_ShouldReturnErrorOnInvalidTime(t *testing.T) {
	// Arrange
	payload := []byte(`{"name":"Bills","description":"Rent","color":"red","date":"2026-07-17","triggers":[{"time":"bad-time","commentary":"Morning reminder"}]}`)

	// Act
	var create CalendarEventCreate
	err := json.Unmarshal(payload, &create)

	// Assert
	require.EqualError(t, err, `invalid time value "bad-time": parsing time "bad-time" as "15:04:05.999999": cannot parse "bad-time" as "15"`)
}

func TestCalendarEventUpdateUnmarshalJSON_ShouldParseSupportedFields(t *testing.T) {
	// Arrange
	triggerID := uuid.New()
	payload := []byte(`{"name":"Bills updated","description":"Rent updated","color":"orange","date":"2026-07-18","triggers":[{"id":"` + triggerID.String() + `","time":"08:45:00","commentary":"Kept and updated"},{"time":"12:00:00","commentary":"Brand new"}]}`)

	// Act
	var update CalendarEventUpdate
	err := json.Unmarshal(payload, &update)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Bills updated", update.Name)
	assert.Equal(t, "Rent updated", update.Description)
	assert.Equal(t, "orange", update.Color)
	assert.True(t, update.Date.Valid)
	assert.Equal(t, "2026-07-18", update.Date.Time.Format("2006-01-02"))
	require.Len(t, update.Triggers, 2)
	require.NotNil(t, update.Triggers[0].Id)
	assert.Equal(t, triggerID, *update.Triggers[0].Id)
	assert.Equal(t, int64(31500000000), update.Triggers[0].Time.Microseconds)
	assert.Nil(t, update.Triggers[1].Id)
	assert.Equal(t, int64(43200000000), update.Triggers[1].Time.Microseconds)
}

func TestCalendarEventUpdateUnmarshalJSON_ShouldReturnErrorOnInvalidTriggerID(t *testing.T) {
	// Arrange
	payload := []byte(`{"name":"Bills updated","description":"Rent updated","color":"orange","date":"2026-07-18","triggers":[{"id":"bad-id","time":"08:45:00","commentary":"Kept and updated"}]}`)

	// Act
	var update CalendarEventUpdate
	err := json.Unmarshal(payload, &update)

	// Assert
	require.EqualError(t, err, "invalid UUID length: 6")
}

func TestCalendarEventDateRangeFilterValidate(t *testing.T) {
	validFrom := newCalendarEventFilterDateValue(2026, time.July, 17)
	validTo := newCalendarEventFilterDateValue(2026, time.July, 18)
	invalidFrom := pgtype.Date{}
	invalidTo := pgtype.Date{}

	tests := []struct {
		name          string
		filter        CalendarEventDateRangeFilter
		expectedError string
	}{
		{
			name: "valid filter",
			filter: CalendarEventDateRangeFilter{
				From: validFrom,
				To:   validTo,
			},
			expectedError: "",
		},
		{
			name: "from is invalid",
			filter: CalendarEventDateRangeFilter{
				From: invalidFrom,
				To:   validTo,
			},
			expectedError: "from should be valid",
		},
		{
			name: "to is invalid",
			filter: CalendarEventDateRangeFilter{
				From: validFrom,
				To:   invalidTo,
			},
			expectedError: "to should be valid",
		},
		{
			name: "from is later than to",
			filter: CalendarEventDateRangeFilter{
				From: validTo,
				To:   validFrom,
			},
			expectedError: "from should not be later than to",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			filter := tt.filter

			// Act
			err := filter.Validate()

			// Assert
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func newCalendarEventFilterDateValue(year int, month time.Month, day int) pgtype.Date {
	return pgtype.Date{
		Time:  time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		Valid: true,
	}
}
