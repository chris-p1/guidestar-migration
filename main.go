package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"

	"golang.org/x/example/hello/reverse"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println(reverse.String("Hello"))

	opt := option.WithCredentialsFile("config/guidestar-multiconfig-firebase-adminsdk-pzpvm-0fcc9de3b9.json")
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://guidestar-multiconfig.firebaseio.com",
	}

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	ref := client.NewRef("38099")
	var data map[string]interface{}

	ref.Get(ctx, &data)
	if err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	fmt.Println(data)

}
