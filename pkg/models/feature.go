package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tschaub/pgfs/pkg/geo"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

// PropertyMap holds arbitrary result properties.
type PropertyMap map[string]interface{}

// Value implements the driver.Valuer interface
func (p PropertyMap) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

// Scan implements the sql.Scanner interface.
func (p *PropertyMap) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) for PropertyMap failed")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("type assertion .(map[string]interface{}) for PropertyMap failed")
	}

	return nil
}

// Feature represents an OGC simple feature.
type Feature struct {
	ID             uuid.UUID    `db:"id"`
	Geometry       geo.Geometry `db:"geometry"`
	Properties     PropertyMap  `db:"properties"`
	CollectionName string       `db:"collection_name"`
}

// Feature implements the Record interface
var _ Record = (*Feature)(nil)

// Features represents a set of Features
type Features []*Feature

// Features implements the RecordSet interface
var _ RecordSet = (*Features)(nil)

// Features implements the BulkInsertable interface
var _ BulkInsertable = (*Features)(nil)

// FeatureQuery is used for querying Features
type FeatureQuery struct {
	Collection Collection
	Limit      uint64
	After      *Feature
}

var defaultFeatureLimit uint64 = 500

// where adds a where clause to the builder based on the query
func (query *FeatureQuery) where(builder sq.SelectBuilder) sq.SelectBuilder {
	builder = builder.
		Where(sq.Eq{column(featureTable, "collection_name"): query.Collection.Name})

	if query.After != nil {
		builder = builder.Where(sq.Gt{column(featureTable, "id"): query.After.ID})
	}

	if query.Limit == 0 {
		query.Limit = defaultFeatureLimit
	}

	return builder.
		OrderBy(fmt.Sprintf("%s ASC", column(featureTable, "id"))).
		Limit(query.Limit + 1)
}

var _ Querier = (*FeatureQuery)(nil)

var featureTable = "features"

var selectFeatures = builder.
	Select(
		column(featureTable, "id"),
		alias(fmt.Sprintf("ST_AsGeoJSON(%s)", column(featureTable, "geometry")), "geometry"),
		column(featureTable, "properties"),
	).
	From(featureTable)

func getFeatureInsertSQL(feature *Feature) (string, []interface{}, error) {
	feature.ID = uuid.New()

	return builder.
		Insert(featureTable).
		SetMap(sq.Eq{
			"id":              feature.ID,
			"collection_name": feature.CollectionName,
			"geometry":        sq.Expr("ST_SetSRID(ST_GeomFromGeoJSON(?), 4326)", feature.Geometry),
			"properties":      feature.Properties,
		}).ToSql()
}

// insert persists a new feature
func (feature *Feature) insert(db *sqlx.DB) error {
	sql, args, sqlErr := getFeatureInsertSQL(feature)
	if sqlErr != nil {
		return sqlErr
	}
	_, execErr := db.Exec(sql, args...)
	return execErr
}

// get retrieves a single feature by ID
func (feature *Feature) get(db *sqlx.DB) error {
	sql, args, err := selectFeatures.Where(sq.Eq{column(featureTable, "id"): feature.ID}).ToSql()
	if err != nil {
		return err
	}
	return db.Get(feature, sql, args...)
}

// update updates a feature's editable fields
func (feature *Feature) update(db *sqlx.DB) error {
	sql, args, sqlErr := builder.
		Update(featureTable).
		SetMap(sq.Eq{
			"geometry":   sq.Expr("ST_SetSRID(ST_GeomFromGeoJSON(?), 4326)", feature.Geometry),
			"properties": feature.Properties,
		}).ToSql()

	if sqlErr != nil {
		return sqlErr
	}

	_, execErr := db.Exec(sql, args...)
	return execErr
}

// delete performs a delete
func (feature *Feature) delete(db *sqlx.DB) error {
	sql, args, sqlErr := builder.
		Delete(featureTable).
		Where(sq.Eq{"id": feature.ID}).ToSql()

	if sqlErr != nil {
		return sqlErr
	}

	_, err := db.Exec(sql, args...)
	return err
}

// insert saves a list of features
func (features *Features) insert(db *sqlx.DB) error {
	tx, txErr := db.Beginx()
	if txErr != nil {
		return txErr
	}

	var err error
	defer func() {
		if err != nil && tx != nil {
			tx.Rollback()
		}
	}()

	for _, feature := range *features {
		sql, args, err := getFeatureInsertSQL(feature)
		if err != nil {
			return err
		}
		_, err = tx.Exec(sql, args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// query gets a list of features
func (features *Features) query(db *sqlx.DB, query Querier) (bool, error) {
	var featureQuery *FeatureQuery
	if query != nil {
		var ok bool
		featureQuery, ok = query.(*FeatureQuery)
		if !ok {
			return false, errors.New("invalid feature query")
		}
	}
	if featureQuery == nil {
		featureQuery = &FeatureQuery{}
	}

	sql, args, err := featureQuery.where(selectFeatures).ToSql()
	if err != nil {
		return false, err
	}

	limit := featureQuery.Limit
	if limit == 0 {
		limit = defaultFeatureLimit
	}

	selectErr := db.Select(features, sql, args...)
	if selectErr != nil {
		return false, selectErr
	}

	more := false
	if len(*features) > int(limit) {
		more = true
		*features = (*features)[:limit]
	}

	return more, nil
}
