package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

var driverName = "postgres"
var builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// Record represents a single database record
type Record interface {
	get(*sqlx.DB) error
	insert(*sqlx.DB) error
	update(*sqlx.DB) error
	delete(*sqlx.DB) error
}

// Querier builds quieries
type Querier interface {
	where(sq.SelectBuilder) sq.SelectBuilder
}

// RecordSet represents a set of database records
type RecordSet interface {
	query(*sqlx.DB, Querier) (bool, error)
}

// BulkInsertable represents a set of records that can be inserted in bulk
type BulkInsertable interface {
	insert(*sqlx.DB) error
}

// Get retrieves a record based on ID
func Get(db *sql.DB, record Record) error {
	return record.get(sqlx.NewDb(db, driverName))
}

// Insert adds a new record and assigns an ID
func Insert(db *sql.DB, record Record) error {
	return record.insert(sqlx.NewDb(db, driverName))
}

// Update sets values for an existing record
func Update(db *sql.DB, record Record) error {
	return record.update(sqlx.NewDb(db, driverName))
}

// Delete removes an existing record
func Delete(db *sql.DB, record Record) error {
	return record.delete(sqlx.NewDb(db, driverName))
}

// Query gets a set of records
func Query(db *sql.DB, records RecordSet, query Querier) (bool, error) {
	return records.query(sqlx.NewDb(db, driverName), query)
}

// BulkInsert inserts a batch of records
func BulkInsert(db *sql.DB, records BulkInsertable) error {
	return records.insert(sqlx.NewDb(db, driverName))
}
