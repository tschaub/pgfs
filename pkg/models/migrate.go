package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var create = `
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS collections (
	name TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS features (
	id UUID PRIMARY KEY,
	collection_name TEXT REFERENCES collections(name) NOT NULL,
	geometry GEOMETRY(GEOMETRY, 4326) NOT NULL,
	properties JSONB NOT NULL
);
CREATE INDEX IF NOT EXISTS features_collection_name_idx ON features(collection_name);
CREATE INDEX IF NOT EXISTS features_geometry_idx ON features USING GIST(geometry);
`

var drop = `
DROP TABLE features;
DROP TABLE collections;
`

// Migrate updates the database
func Migrate(db *sql.DB) error {
	sqlxDB := sqlx.NewDb(db, driverName)
	_, err := sqlxDB.Exec(create)

	return err
}

// Drop all the data
func Drop(db *sql.DB) error {
	sqlxDB := sqlx.NewDb(db, driverName)
	_, err := sqlxDB.Exec(drop)
	return err
}
