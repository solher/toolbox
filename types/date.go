package types

import (
	"database/sql/driver"
	"errors"
	"time"
)

// Date is a string which transparently converts itself to a sql date.
type Date string

const dateFormat = "2006-01-02"

// Value implements the driver.Valuer interface, converting the Date to a string.
func (d Date) Value() (driver.Value, error) {
	if d == "" {
		return nil, nil
	}
	return string(d), nil
}

// Scan implements the sql.Scanner interface, converting a sql date to a Date.
func (d *Date) Scan(value interface{}) error {
	switch t := value.(type) {
	case []uint8:
		*d = Date(t)
	case time.Time:
		*d = Date(t.Format(dateFormat))
	case nil:
	default:
		return errors.New("Incompatible type for Date")
	}

	return nil
}
