package actions

import (
	"fmt"

	"github.com/gobuffalo/buffalo/worker"
	pgTypes "github.com/mc2soft/pq-types"
	"github.com/wung-s/gotv/models"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var w worker.Worker

func init() {
	w = worker.NewSimple()
	w.Register("geocode_address", func(args worker.Args) error {
		memberID := fmt.Sprint(args["memberID"])
		address := fmt.Sprint(args["address"])

		tx := models.DB
		member := &models.Member{}
		if err := tx.Find(member, memberID); err != nil {
			return err
		}

		gr := &maps.GeocodingRequest{
			Address: address,
		}

		resp, err := gMap.Geocode(context.Background(), gr)
		if err != nil {
			fmt.Print("Error geocoding", err)
		}

		coords := resp[0].Geometry.Location
		member.LatLng = pgTypes.PostGISPoint{Lon: coords.Lng, Lat: coords.Lat}
		_, err = tx.ValidateAndUpdate(member)
		return err
	})
}
