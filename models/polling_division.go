package models

import (
	"github.com/gobuffalo/uuid"
)

// PollingDivision reflects the db structure
type PollingDivision struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Edid        int       `json:"edid" db:"edid"`
	No          int       `json:"no" db:"no"`
	ShapeArea   string    `json:"shape_area" db:"shape_area"`
	ShapeLength string    `json:"shape_length" db:"shape_length"`
}

// PollingDivisions is not required by pop and may be deleted
type PollingDivisions []PollingDivision
