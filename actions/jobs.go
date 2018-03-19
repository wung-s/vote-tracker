package actions

import (
	"encoding/json"
	"fmt"

	"github.com/gobuffalo/buffalo/worker"
	uuid "github.com/gobuffalo/uuid"
	pgTypes "github.com/mc2soft/pq-types"
	"github.com/wung-s/gotv/models"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var w worker.Worker

func init() {
	w = worker.NewSimple()
	w.Register("geocode_address", func(args worker.Args) error {
		var memberIDs []uuid.UUID
		if err := json.Unmarshal([]byte(fmt.Sprint(args["memberIds"])), &memberIDs); err != nil {
			fmt.Println("Error obtaining memberIDs", err)
			return err
		}

		tx := models.DB
		// obtain geolocation for each newly created members
		for _, id := range memberIDs {
			member := &models.Member{}
			if err := tx.Find(member, id); err != nil {
				return err
			}
			gr := &maps.GeocodingRequest{
				Address: member.Address(),
			}

			resp, err := gMap.Geocode(context.Background(), gr)
			if err != nil {
				fmt.Print("Error geocoding", err)
				return err
			}

			coords := resp[0].Geometry.Location
			member.LatLng = pgTypes.PostGISPoint{Lon: coords.Lng, Lat: coords.Lat}
			if _, err = tx.ValidateAndUpdate(member); err != nil {
				return err
			}
		}
		return nil
	})
}
