package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/lib/pq"
	"github.com/tschaub/pgfs/pkg/models"
)

// CollectionInfo encodes collection information
type CollectionInfo struct {
	Name        string `json:"name" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// CollectionList encodes a list of collections
type CollectionList struct {
	Collections []*CollectionInfo `json:"collections"`
}

// CreateCollection saves a new collection
func CreateCollection(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		info := &CollectionInfo{}
		if bindErr := c.Bind(info); bindErr != nil {
			return bindErr
		}

		if validateErr := c.Validate(info); validateErr != nil {
			return validateErr
		}

		collection := &models.Collection{
			Name:        info.Name,
			Title:       info.Title,
			Description: info.Description,
		}

		createErr := models.Insert(db, collection)
		if createErr != nil {
			if pqErr, ok := createErr.(*pq.Error); ok {
				if pqErr.Code.Name() == "unique_violation" {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Collection with name '%s' already exists", info.Name))
				}
			}
			return createErr
		}

		return c.JSON(http.StatusCreated, info)
	}
}

// GetCollection responds with a single collection by name
func GetCollection(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")

		collection := &models.Collection{Name: name}
		getErr := models.Get(db, collection)
		if getErr != nil {
			if getErr == sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusNotFound)
			}
			return getErr
		}

		info := &CollectionInfo{
			Name:        collection.Name,
			Title:       collection.Title,
			Description: collection.Description,
		}

		return c.JSON(http.StatusOK, info)
	}
}

// ListCollections responds with a list of all the collections
func ListCollections(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		collections := models.Collections{}
		_, listErr := models.Query(db, &collections, nil)
		if listErr != nil {
			return listErr
		}

		list := make([]*CollectionInfo, len(collections))
		for i, collection := range collections {
			list[i] = &CollectionInfo{
				Name:        collection.Name,
				Title:       collection.Title,
				Description: collection.Description,
			}
		}

		return c.JSON(http.StatusOK, &CollectionList{Collections: list})
	}
}
