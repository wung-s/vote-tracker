package actions

import (
	"errors"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// FirebaseApp holds the connection to Firebase, used for authentication
var FirebaseApp *firebase.App

// InitializeFirebase set up firebase
func InitializeFirebase() error {
	err := errors.New("Placeholder error")

	// opt := option.WithCredentialsFile(dir + "/serviceAccountKey.json")
	opt := option.WithCredentialsFile(os.Getenv("FB_SERVICE_AC_KEY"))
	FirebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return err
}
