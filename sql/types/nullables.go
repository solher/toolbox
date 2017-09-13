package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

// NullFloat64 is a float64 that is NULL when set to its zero.
type NullFloat64 float64

// Value implements the driver.Valuer interface.
func (n NullFloat64) Value() (driver.Value, error) {
	if n == 0 {
		return nil, nil
	}
	return float64(n), nil
}

// Scan implements the sql.Scanner interface.
func (n *NullFloat64) Scan(value interface{}) error {
	switch value := value.(type) {
	case nil:
		*n = 0
	case string:
		tmp, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		*n = NullFloat64(tmp)
	case []byte:
		tmp, err := strconv.ParseFloat(string(value), 64)
		if err != nil {
			return err
		}
		*n = NullFloat64(tmp)
	case float64:
		*n = NullFloat64(value)
	case float32:
		*n = NullFloat64(value)
	default:
		return fmt.Errorf("Incompatible type for NullFloat64")
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
		tmp, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*n = NullInt64(tmp)
	case []byte:
		tmp, err := strconv.ParseInt(string(value), 10, 64)
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

// NullUint64 is an unsigned int that is NULL when set to its zero.
type NullUint64 uint64

// Value implements the driver.Valuer interface.
func (n NullUint64) Value() (driver.Value, error) {
	if n == 0 {
		return nil, nil
	}
	return int64(n), nil
}

// Scan implements the sql.Scanner interface.
func (n *NullUint64) Scan(value interface{}) error {
	switch value := value.(type) {
	case nil:
		*n = 0
	case string:
		tmp, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*n = NullUint64(tmp)
	case []byte:
		tmp, err := strconv.ParseUint(string(value), 10, 64)
		if err != nil {
			return err
		}
		*n = NullUint64(tmp)
	case int64:
		*n = NullUint64(value)
	case int32:
		*n = NullUint64(value)
	case uint64:
		*n = NullUint64(value)
	default:
		return fmt.Errorf("Incompatible type for NullUint64")
	}
	return nil
}

// NullTimestamp is a time that is NULL when set to its zero.
type NullTimestamp struct {
	time.Time
}

// Value implements the driver.Valuer interface.
func (n NullTimestamp) Value() (driver.Value, error) {
	if n.IsZero() {
		return nil, nil
	}
	return n.Time, nil
}

// Scan implements the sql.Scanner interface.
func (n *NullTimestamp) Scan(value interface{}) error {
	switch value := value.(type) {
	case nil:
		*n = NullTimestamp{Time: time.Time{}}
	case time.Time:
		*n = NullTimestamp{Time: value}
	default:
		return fmt.Errorf("Incompatible type for NullTimestamp")
	}
	return nil
}

// NullTime is a string formatted as a Postgres time without timestamp that is NULL when set to its zero.
type NullTime string

const timeFormat = "15:04:05"

// Value implements the driver.Valuer interface.
func (n NullTime) Value() (driver.Value, error) {
	if n == "" {
		return nil, nil
	}
	return string(n), nil
}

// Scan implements the sql.Scanner interface.
func (n *NullTime) Scan(value interface{}) error {
	switch t := value.(type) {
	case nil:
		*n = ""
	case []uint8:
		*n = NullTime(t)
	case time.Time:
		*n = NullTime(t.Format(timeFormat))
	default:
		return errors.New("Incompatible type for NullTime")
	}
	return nil
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
	case nil:
		*n = ""
	case []uint8:
		*n = NullDate(t)
	case time.Time:
		*n = NullDate(t.Format(dateFormat))
	default:
		return errors.New("Incompatible type for NullDate")
	}
	return nil
}

// NullMoney is a float64 formatted as a Postgres money type that is NULL when set to its zero.
type NullMoney float64

// Value implements the driver.Valuer interface.
func (n NullMoney) Value() (driver.Value, error) {
	if n == 0.0 {
		return nil, nil
	}
	return float64(n), nil
}

// Scan implements the sql.Scanner interface.
func (n *NullMoney) Scan(value interface{}) error {
	switch t := value.(type) {
	case nil:
		*n = 0.0
	case string:
		money := strings.Replace(t, ",", "", -1)
		money = strings.Replace(money, "$", "", -1)
		m, err := strconv.ParseFloat(money, 64)
		if err != nil {
			return err
		}
		*n = NullMoney(m)
	case []uint8:
		money := strings.Replace(string(t), ",", "", -1)
		money = strings.Replace(money, "$", "", -1)
		m, err := strconv.ParseFloat(money, 64)
		if err != nil {
			return err
		}
		*n = NullMoney(m)
	default:
		return errors.New("Incompatible type for NullMoney")
	}
	return nil
}
