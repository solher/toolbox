package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// NullString is a string that is NULL when set to its zero.
type NullString string

// Value implements the driver.Valuer interface.
func (n NullString) Value() (driver.Value, error) {
	if n == "" {
		return nil, nil
	}
	return string(n), nil
}

// Scan implements the sql.Scanner interface.
func (n *NullString) Scan(value interface{}) error {
	switch value := value.(type) {
	case nil:
		*n = ""
	case string:
		*n = NullString(value)
	case []byte:
		*n = NullString(value)
	default:
		return fmt.Errorf("Incompatible type for NullString")
	}
	return nil
}

// NullInt64 is an int that is NULL when set to its zero.
type NullInt64 int64

// Value implements the driver.Valuer interface.
func (n NullInt64) Value() (driver.Value, error) {
	if n == 0 {
		return nil, nil
	}
	return int64(n), nil
}

// Scan implements the sql.Scanner interface.
func (n *NullInt64) Scan(value interface{}) error {
	switch value := value.(type) {
	case nil:
		*n = 0
	case string:
		tmp, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		*n = NullInt64(tmp)
	case []byte:
		tmp, err := strconv.ParseUint(string(value), 10, 32)
		if err != nil {
			return err
		}
		*n = NullInt64(tmp)
	case int64:
		*n = NullInt64(value)
	case int32:
		*n = NullInt64(value)
	case uint64:
		*n = NullInt64(value)
	default:
		return fmt.Errorf("Incompatible type for NullInt64")
	}
	return nil
}

// NullTime is a time that is NULL when set to its zero.
type NullTime time.Time

// Value implements the driver.Valuer interface.
func (n NullTime) Value() (driver.Value, error) {
	if time.Time(n).IsZero() {
		return nil, nil
	}
	return time.Time(n), nil
}

// Scan implements the sql.Scanner interface.
func (n *NullTime) Scan(value interface{}) error {
	switch value := value.(type) {
	case nil:
		*n = NullTime(time.Time{})
	case time.Time:
		*n = NullTime(value)
	default:
		return fmt.Errorf("Incompatible type for NullTime")
	}
	return nil
}

// IsZero is there to facilitate usage with Go templates.
func (n NullTime) IsZero() bool {
	return time.Time(n).IsZero()
}

// MarshalJSON implements the json.Marshaler interface.
func (n NullTime) MarshalJSON() ([]byte, error) {
	return time.Time(n).MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *NullTime) UnmarshalJSON(data []byte) error {
	return (*time.Time)(n).UnmarshalJSON(data)
}

// NullDate is a string formatted as a Postgres date that is NULL when set to its zero.
type NullDate string

const dateFormat = "2006-01-02"

// Value implements the driver.Valuer interface.
func (n NullDate) Value() (driver.Value, error) {
	if n == "" {
		return nil, nil
	}
	return string(n), nil
}

// Scan implements the sql.Scanner interface.
func (n *NullDate) Scan(value interface{}) error {
	switch t := value.(type) {
	case []uint8:
		*n = NullDate(t)
	case time.Time:
		*n = NullDate(t.Format(dateFormat))
	case nil:
	default:
		return errors.New("Incompatible type for NullDate")
	}

	return nil
}
