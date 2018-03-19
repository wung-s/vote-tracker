package models

import (
	"github.com/gobuffalo/uuid"
	pgTypes "github.com/mc2soft/pq-types"
)

// ElectoralDistrict reflects the db structure
type ElectoralDistrict struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Edid        int                    `json:"edid" db:"edid"`
	ShapeArea   string                 `json:"shape_area" db:"shape_area"`
	ShapeLength string                 `json:"shape_length" db:"shape_length"`
	Geom        pgTypes.PostGISPolygon `json:"geom" db:"geom"`
}

// ElectoralDistricts is not required by pop and may be deleted
type ElectoralDistricts []ElectoralDistrict
