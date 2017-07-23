package sql

import "github.com/jmoiron/sqlx"

// Transaction provides a simple API handling commits and rollbacks.
func Transaction(db *sqlx.DB, transaction func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	if err := transaction(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}
