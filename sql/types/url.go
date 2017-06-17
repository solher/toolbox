package types

import (
	"database/sql/driver"
	"errors"
	"net/url"
)

// URL is a url.URL which transparently converts itself to a string or byte array.
type URL url.URL

// Value implements the driver.Valuer interface, converting the URL to a string.
func (u URL) Value() (driver.Value, error) {
	val := url.URL(u)
	return val.String(), nil
}

// Scan implements the sql.Scanner interface, converting a string or byte array to a URL.
func (u *URL) Scan(src interface{}) error {
	var rawurl string
	switch src := src.(type) {
	case string:
		rawurl = src
	case []byte:
		rawurl = string(src)
	default:
		return errors.New("Incompatible type for URL")
	}

	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	*u = URL(*parsedURL)
	return nil
}
