package actions

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gobuffalo/packr"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// FirebaseApp holds the connection to Firebase, used for authentication
var FirebaseApp *firebase.App

// InitializeFirebase set up firebase
func InitializeFirebase() error {
	err := errors.New("Placeholder error")

	box := packr.NewBox("../config")

	content := []byte{}
	if ENV == "test" {
		content = box.Bytes("serviceAccountKey.json")
	} else {
		content = box.Bytes(os.Getenv("FB_SERVICE_AC_KEY"))
	}

	fileName := "firebaseKey.json"
	if err := ioutil.WriteFile(fileName, content, 0644); err != nil {
		fmt.Println("Error writing firebasekey file:", err)
	}

	opt := option.WithCredentialsFile(fileName)

	FirebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return err
}
