package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Time is a string formatted as a Postgres time without timestamp that is NULL when set to its zero.
type Time struct {
	time.Time
}

// String returns the time as a string formatted as a time.
func (n Time) String() string {
	return n.Time.Format(time.TimeOnly)
}

// MarshalJSON marshals the time as a string formatted as a time.
func (n Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

// UnmarshalJSON unmarshals the time as a string formatted as a time.
func (n *Time) UnmarshalJSON(data []byte) error {
	var t string
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	parsed, err := time.Parse(time.TimeOnly, t)
	if err != nil {
		// Just in case, we also try to parse hh:mm
		if parsed, err := time.Parse("15:04", t); err == nil {
			*n = Time{Time: parsed}
			return nil
		}
		return err
	}
	*n = Time{Time: parsed}
	return nil
}

// Value implements the driver.Valuer interface.
func (n Time) Value() (driver.Value, error) {
	return n.String(), nil
}

// Scan implements the sql.Scanner interface.
func (n *Time) Scan(value any) error {
	switch t := value.(type) {
	case nil:
		*n = Time{Time: time.Time{}}
	case string:
		parsed, err := time.Parse(time.TimeOnly, t)
		if err != nil {
			// Just in case, we also try to parse hh:mm
			if parsed, err := time.Parse("15:04", t); err == nil {
				*n = Time{Time: parsed}
				return nil
			}
			return err
		}
		*n = Time{Time: parsed}
	default:
		return errors.New("incompatible type for Time")
	}
	return nil
}

// Date is a string formatted as a Postgres date that is NULL when set to its zero.
type Date struct {
	time.Time
}

// String returns the date as a string formatted as a date.
func (n Date) String() string {
	return n.Time.Format(time.DateOnly)
}

// MarshalJSON marshals the date as a string formatted as a date.
func (n Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

// UnmarshalJSON unmarshals the date as a string formatted as a date.
func (n *Date) UnmarshalJSON(data []byte) error {
	var t string
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	parsed, err := time.Parse(time.DateOnly, t)
	if err != nil {
		return err
	}
	*n = Date{Time: parsed}
	return nil
}

// Value implements the driver.Valuer interface.
func (n Date) Value() (driver.Value, error) {
	return n.String(), nil
}

// Scan implements the sql.Scanner interface.
func (n *Date) Scan(value any) error {
	switch t := value.(type) {
	case nil:
		*n = Date{Time: time.Time{}}
	case string:
		parsed, err := time.Parse(time.DateOnly, t)
		if err != nil {
			return err
		}
		*n = Date{Time: parsed}
	default:
		return errors.New("incompatible type for Date")
	}
	return nil
}

// TimeZone is a string formatted as a Postgres time zone that is NULL when set to its zero.
type TimeZone struct {
	time.Location
}

// String returns the time zone as a string formatted as a time zone.
func (n TimeZone) String() string {
	return n.Location.String()
}

// MarshalJSON marshals the time zone as a string formatted as a time zone.
func (n TimeZone) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

// UnmarshalJSON unmarshals the time zone as a string formatted as a time zone.
func (n *TimeZone) UnmarshalJSON(data []byte) error {
	var t string
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	parsed, err := time.LoadLocation(t)
	if err != nil {
		return err
	}
	*n = TimeZone{Location: *parsed}
	return nil
}

// Value implements the driver.Valuer interface.
func (n TimeZone) Value() (driver.Value, error) {
	return n.Location.String(), nil
}

// Scan implements the sql.Scanner interface.
func (n *TimeZone) Scan(value any) error {
	switch t := value.(type) {
	case nil:
		*n = TimeZone{Location: *time.UTC}
	case string:
		parsed, err := time.LoadLocation(t)
		if err != nil {
			return err
		}
		*n = TimeZone{Location: *parsed}
	default:
		return errors.New("incompatible type for TimeZone")
	}
	return nil
}
