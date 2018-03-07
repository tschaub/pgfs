package geo

import (
	"database/sql/driver"
	"encoding/json"

	geojson "github.com/paulmach/go.geojson"
)

// Geometry represents a GeoJSON geometry
type Geometry struct {
	geometry geojson.Geometry
}

// Valid determines if a geometry is valid
func (g *Geometry) Valid() bool {
	_, err := g.MarshalJSON()
	return err == nil
}

// UnmarshalJSON decodes the data into a geometry
func (g *Geometry) UnmarshalJSON(data []byte) error {
	var geometry geojson.Geometry
	err := json.Unmarshal(data, &geometry)
	if err != nil {
		return err
	}
	g.geometry = geometry

	return nil
}

// MarshalJSON encodes the geometry into GeoJSON
func (g *Geometry) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.geometry)
}

// Scan implements the sql.Scanner interface
func (g *Geometry) Scan(value interface{}) error {
	return g.geometry.Scan(value)
}

// Value implements the driver.Valuer interface
func (g Geometry) Value() (driver.Value, error) {
	data, err := g.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return string(data), nil
}
