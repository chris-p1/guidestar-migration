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

type configId string

type wsConfig struct {
	Settings map[string]Config `json:"settings"`
}

type Config struct {
	LoggingField509A           int    `json:"loggingField509a"`
	LoggingFieldEin            int    `json:"loggingFieldEin"`
	LoggingFieldLookupDate     int    `json:"loggingFieldLookupDate"`
	LoggingFieldPub78          int    `json:"loggingFieldPub78"`
	LoggingFieldRulingDate     int    `json:"loggingFieldRulingDate"`
	LoggingFieldSubsectionDesc int    `json:"loggingFieldSubsectionDesc"`
	LoggingFormID              int    `json:"loggingFormId"`
	Name                       string `json:"name"`
	TargetFieldID              int    `json:"targetFieldId"`
	TargetFormID               int    `json:"targetFormId"`
}

func migrate(config wsConfig) wsConfig {
	// var newConfig map[string]interface{}
	settings := config.Settings
	for configId, configVals := range settings {
		fmt.Printf("Key: %s, Value: %s\n", configId, configVals.Name)
	}
	return config
}

func main() {
	fmt.Println(reverse.String("Guidestar Migration"))

	// 20250714T114955 TODO: change to prod once ready
	testWs := "33940"
	// testWs2 := "40138"
	// testWsProd := "38099"
	fbUrl := "https://guidestar-multiconfig.firebaseio.com"
	opt := option.WithCredentialsFile("config/guidestar-multiconfig-firebase-adminsdk-pzpvm-0fcc9de3b9.json")
	// fbUrl := "https://guidestar-stage.firebaseio.com"
	// opt := option.WithCredentialsFile("config/guidestar-stage-firebase-adminsdk-7j26g-c38da55334.json")
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

	var oldConfig wsConfig
	oldRef := client.NewRef(testWs)
	oldRef.Get(ctx, &oldConfig)
	if err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	newConfig := migrate(oldConfig)
	// updatedRef := client.NewRef(testWs2)
	// var updatedConfig map[string]interface{}

	// updatedRef.Get(ctx, &updatedConfig)
	// if err != nil {
	// 	log.Fatalln("Error reading from database:", err)
	// }

	oldTest, _ := json.MarshalIndent(newConfig, "", "  ")
	// updatedTest, _ := json.MarshalIndent(updatedConfig, "", "  ")
	fmt.Println(string(oldTest))
	// fmt.Println(string(updatedTest))

}
