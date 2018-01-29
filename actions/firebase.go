package actions

import (
	"fmt"
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
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	opt := option.WithCredentialsFile(dir + "/serviceAccountKey.json")
	FirebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return err
}
