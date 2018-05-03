package actions

import (
	"errors"
	"log"
	"os"

	"googlemaps.github.io/maps"
)

var gMap *maps.Client

// InitializeGoogleMaps instantiates an app-wide Google maps client
func InitializeGoogleMaps() error {
	err := errors.New("")
	gMap, err = maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_KEY")))
	if err != nil {
		log.Fatalf("error initializing Google Maps: %v\n", err)
	}
	return err
}
