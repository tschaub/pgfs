package handlers

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/tschaub/pgfs/pkg/geo"
	"github.com/tschaub/pgfs/pkg/models"
)

// NewFeatureInfo represents a GeoJSON Feature
type NewFeatureInfo struct {
	Geometry   geo.Geometry           `json:"geometry" validate:"required"`
	Properties map[string]interface{} `json:"properties" validate:"required"`
}

// NewFeatureList is a GeoJSON FeatureCollection
type NewFeatureList struct {
	Type     string            `json:"type" validate:"required"`
	Features []*NewFeatureInfo `json:"features" validate:"required"`
	More     bool              `json:"more"`
}

// FeatureInfo represents a GeoJSON Feature
type FeatureInfo struct {
	ID         uuid.UUID              `json:"id"`
	Geometry   geo.Geometry           `json:"geometry" validate:"required"`
	Properties map[string]interface{} `json:"properties" validate:"required"`
}

// FeatureList is a GeoJSON FeatureCollection
type FeatureList struct {
	Type     string         `json:"type" validate:"required"`
	Features []*FeatureInfo `json:"features" validate:"required"`
	More     bool           `json:"more"`
}

// FeatureListQuery allows features to be queried
type FeatureListQuery struct {
	Count uint64 `query:"count"`
	After string `query:"after"`
}

func infoFromFeature(f *models.Feature) *FeatureInfo {
	return &FeatureInfo{
		ID:         f.ID,
		Geometry:   f.Geometry,
		Properties: f.Properties,
	}
}

// ListFeatures responds with a list of features
func ListFeatures(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		query := &FeatureListQuery{}
		bindErr := c.Bind(query)
		if bindErr != nil {
			return bindErr
		}

		name := c.Param("collectionName")
		collection := &models.Collection{Name: name}
		if getErr := models.Get(db, collection); getErr != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		featureQuery := &models.FeatureQuery{
			Collection: *collection,
			Limit:      query.Count,
		}

		if query.After != "" {
			id, err := uuid.Parse(query.After)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "cannot parse 'after' as a UUID")
			}
			feature := &models.Feature{ID: id}
			getErr := models.Get(db, feature)
			if getErr != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "bad 'after' id")
			}
			featureQuery.After = feature
		}

		features := models.Features{}
		more, listErr := models.Query(db, &features, featureQuery)
		if listErr != nil {
			return listErr
		}

		list := make([]*FeatureInfo, len(features))
		for i, result := range features {
			list[i] = infoFromFeature(result)
		}

		resultList := &FeatureList{
			Type:     "FeatureCollection",
			Features: list,
			More:     more,
		}

		return c.JSON(http.StatusOK, resultList)
	}
}

// AddFeatures adds features to a collection
func AddFeatures(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("collectionName")

		collection := &models.Collection{Name: name}
		getErr := models.Get(db, collection)
		if getErr != nil {
			if getErr == sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusNotFound)
			}
			return getErr
		}

		if collection.Name != name {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		info := &NewFeatureList{}
		if bindErr := c.Bind(info); bindErr != nil {
			return bindErr
		}

		features := make(models.Features, len(info.Features))
		for i, feature := range info.Features {
			features[i] = &models.Feature{
				ID:             uuid.New(),
				CollectionName: name,
				Geometry:       feature.Geometry,
				Properties:     feature.Properties,
			}
		}

		return models.BulkInsert(db, &features)
	}
}
