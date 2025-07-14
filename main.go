package main

/*
What is this?

This is a script to migrate the Firebase data from the old Guidestar plugin to
the new one. The old one has a slightly older format that needs to be tweaked a
bit to work with the new plugin.
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"

	"golang.org/x/example/hello/reverse"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println(reverse.String("Guidestar Migration"))

	// 20250714T114955 TODO: change to prod once ready
	testWs := "15032"
	// fbUrl := "https://guidestar-multiconfig.firebaseio.com"
	// opt := option.WithCredentialsFile("config/guidestar-multiconfig-firebase-adminsdk-pzpvm-0fcc9de3b9.json")
	fbUrl := "https://guidestar-stage.firebaseio.com"
	opt := option.WithCredentialsFile("config/guidestar-stage-firebase-adminsdk-7j26g-c38da55334.json")
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: fbUrl,
	}

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	ref := client.NewRef(testWs)
	var data map[string]interface{}

	ref.Get(ctx, &data)
	if err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	test, _ := json.Marshal(data)
	fmt.Println(string(test))

}
