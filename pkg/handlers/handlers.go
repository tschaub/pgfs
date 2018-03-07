package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	validator "gopkg.in/go-playground/validator.v9"
)

// Validator provides request body validation
type Validator struct {
	validator *validator.Validate
}

// Validate checks structs based on validate tags
func (v *Validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err != nil {
		err = echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return err
}

// New creates a new handler
func New(db *sql.DB) *echo.Echo {
	router := echo.New()
	router.HideBanner = true

	router.Validator = &Validator{validator: validator.New()}

	// 500 on panic
	router.Use(middleware.Recover())

	// make smaller
	router.Use(middleware.Gzip())

	// set up cors
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderContentLength, echo.HeaderAuthorization},
		MaxAge:       24 * 60 * 60,
	}))

	// list collections
	router.GET("/collections", ListCollections(db))

	// create new collection
	router.POST("/collections", CreateCollection(db))

	// get a single collection
	router.GET("/collections/:name", GetCollection(db))

	// add features to collection
	router.POST("/collections/:collectionName/items", AddFeatures(db))

	// list features for a collection
	router.GET("/collections/:collectionName/items", ListFeatures(db))

	return router
}
