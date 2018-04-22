package actions

import (
	"encoding/json"
	"fmt"

	"github.com/gobuffalo/buffalo/worker"
	"github.com/gobuffalo/pop/nulls"
	uuid "github.com/gobuffalo/uuid"
	pgTypes "github.com/mc2soft/pq-types"
	"github.com/pkg/errors"
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
			fmt.Println(err)
			return err
		}

		tx := models.DB
		// obtain geolocation for each newly created members
		for _, id := range memberIDs {
			member := &models.Member{}
			if err := tx.Find(member, id); err != nil {
				fmt.Printf("Error finding member with ID: %v %v", id, err)
				continue
			}
			gr := &maps.GeocodingRequest{
				Address: member.Address(),
			}

			resp, err := gMap.Geocode(context.Background(), gr)
			if err != nil {
				fmt.Print(err)
				return err
			}
			// check for if Google Maps API returned any geocode data for the supplied address
			if len(resp) > 0 {
				coords := resp[0].Geometry.Location
				sql := `select id from polling_divisions
     					where ST_Contains(polling_divisions.geom, ST_GeomFromText(?,4326))`
				point := fmt.Sprintf("POINT(%v %v)", coords.Lng, coords.Lat)

				result := struct {
					ID nulls.UUID
				}{}
				var exist bool
				if exist, err = tx.RawQuery(sql, point).Exists(&result); err != nil {
					fmt.Println(err)
					return err
				}

				if exist {
					tx.RawQuery(sql, point).First(&result)
					member.PollingDivisionID = result.ID
				}

				member.LatLng = pgTypes.PostGISPoint{Lon: coords.Lng, Lat: coords.Lat}
				if _, err = tx.ValidateAndUpdate(member); err != nil {
					return err
				}
			} else {
				errors.New("No geocode information returned")
			}
		}
		return nil
	})
}
