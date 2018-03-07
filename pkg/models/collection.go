package models

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

// Collection is a set of features.
type Collection struct {
	Name        string `db:"name"`
	Title       string `db:"title"`
	Description string `db:"description"`
}

// Collection implements the Record interface
var _ Record = (*Collection)(nil)

// Collections represents a set of collections
type Collections []*Collection

// CollectionsQuery is used for querying collections
type CollectionsQuery struct {
	Limit uint
}

// where adds a where clause to the builder based on the query
func (query *CollectionsQuery) where(builder sq.SelectBuilder) sq.SelectBuilder {
	// TODO: query collections
	return builder
}

// CollectionsQuery implements the Querier interface
var _ Querier = (*CollectionsQuery)(nil)

// Collections implements the RecordSet interface
var _ RecordSet = (*Collections)(nil)

var collectionTable = "collections"

var selectCollections = builder.
	Select(
		column(collectionTable, "name"),
		column(collectionTable, "title"),
		column(collectionTable, "description")).
	From(collectionTable).
	OrderBy(fmt.Sprintf("%s ASC", column(collectionTable, "name")))

// insert persists a new collection
func (collection *Collection) insert(db *sqlx.DB) error {
	sql, args, sqlErr := builder.
		Insert(collectionTable).
		SetMap(sq.Eq{
			"title":       collection.Title,
			"name":        collection.Name,
			"description": collection.Description,
		}).ToSql()

	if sqlErr != nil {
		return sqlErr
	}

	_, insertErr := db.Exec(sql, args...)
	if insertErr != nil {
		return insertErr
	}

	return collection.get(db)
}

// update updates a collection's editable fields
func (collection *Collection) update(db *sqlx.DB) error {
	sql, args, sqlErr := builder.
		Update(collectionTable).
		SetMap(sq.Eq{
			"title":       collection.Title,
			"description": collection.Description,
		}).
		Where(sq.Eq{"name": collection.Name}).ToSql()

	if sqlErr != nil {
		return sqlErr
	}

	_, execErr := db.Exec(sql, args...)
	return execErr
}

// get finds a collection by name
func (collection *Collection) get(db *sqlx.DB) error {
	sql, args, sqlErr := selectCollections.Where(sq.Eq{column(collectionTable, "name"): collection.Name}).ToSql()
	if sqlErr != nil {
		return sqlErr
	}

	return db.Get(collection, sql, args...)
}

// delete removes the collection
func (collection *Collection) delete(db *sqlx.DB) error {
	// TODO: delete features first

	sql, args, sqlErr := builder.
		Delete(collectionTable).
		Where(sq.Eq{"name": collection.Name}).ToSql()

	if sqlErr != nil {
		return sqlErr
	}

	_, err := db.Exec(sql, args...)
	return err
}

// query lists collections that match the given query (or all if nil)
func (collections *Collections) query(db *sqlx.DB, query Querier) (bool, error) {
	var collectionQuery *CollectionsQuery
	if query != nil {
		var ok bool
		collectionQuery, ok = query.(*CollectionsQuery)
		if !ok {
			return false, errors.New("invalid project query")
		}
	}
	if collectionQuery == nil {
		collectionQuery = &CollectionsQuery{}
	}

	sql, args, err := collectionQuery.where(selectCollections).ToSql()
	if err != nil {
		return false, err
	}

	selectErr := db.Select(collections, sql, args...)
	if selectErr != nil {
		return false, selectErr
	}

	return false, nil
}
