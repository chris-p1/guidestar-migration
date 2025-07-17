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

type NewFormFields struct {
	LoggingField509A           string `json:"loggingField509a"`
	LoggingFieldLookupDate     string `json:"loggingFieldLookupDate"`
	LoggingFieldPub78          string `json:"loggingFieldPub78"`
	LoggingFieldRulingDate     string `json:"loggingFieldRulingDate"`
	LoggingFieldSubsectionDesc string `json:"loggingFieldSubsectionDesc"`
}

type NewConfigFields struct {
	Enabled                    bool   `json:"enabled"`
	LoggingField509A           string `json:"loggingField509a,omitempty"`
	LoggingFieldLookupDate     string `json:"loggingFieldLookupDate,omitempty"`
	LoggingFieldPub78          string `json:"loggingFieldPub78,omitempty"`
	LoggingFieldRulingDate     string `json:"loggingFieldRulingDate,omitempty"`
	LoggingFieldSubsectionDesc string `json:"loggingFieldSubsectionDesc,omitempty"`
	LoggingFormEnabled         bool   `json:"loggingFormEnabled,omitempty"`
	LoggingFormID              int    `json:"loggingFormId,omitempty"`
	LoggingLinkedFormID        string `json:"loggingLinkedFormId,omitempty"`
	Mch1                       Mch1   `json:"mch1"`
	Name                       string `json:"name"`
	TargetFieldID              int    `json:"targetFieldId"`
	TargetFormID               int    `json:"targetFormId"`
}

type Mch1 struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type OldConfigFields struct {
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

type OldConfig struct {
	Settings map[string]OldConfigFields `json:"settings"`
}

type NewConfig struct {
	Forms    map[string]NewFormFields   `json:"forms"`
	Settings map[string]NewConfigFields `json:"settings"`
}

func migrate(oldConfig OldConfig) (NewConfig, error) {
	// var newConfig map[string]interface{}
	newConf := NewConfig{
		Forms:    make(map[string]NewFormFields),
		Settings: make(map[string]NewConfigFields),
	}
	oldSettings := oldConfig.Settings

	for configId := range oldSettings {

		fmt.Printf("Key: %s\nValue: %#v\n", configId, oldSettings[configId])

		newConf.Settings[configId] = NewConfigFields{
			Enabled:                    true,
			LoggingField509A:           fmt.Sprint(oldSettings[configId].LoggingField509A),
			LoggingFieldLookupDate:     fmt.Sprint(oldSettings[configId].LoggingFieldLookupDate),
			LoggingFieldPub78:          fmt.Sprint(oldSettings[configId].LoggingFieldPub78),
			LoggingFieldRulingDate:     fmt.Sprint(oldSettings[configId].LoggingFieldRulingDate),
			LoggingFieldSubsectionDesc: fmt.Sprint(oldSettings[configId].LoggingFieldSubsectionDesc),
			LoggingFormEnabled:         true,
			LoggingFormID:              oldSettings[configId].LoggingFormID,
			LoggingLinkedFormID:        fmt.Sprint(oldSettings[configId].LoggingFieldEin),
			Mch1:                       Mch1{Type: "Logging", Value: "Yes"},
			Name:                       oldSettings[configId].Name,
			TargetFieldID:              oldSettings[configId].TargetFieldID,
			TargetFormID:               oldSettings[configId].TargetFormID,
		}

		logFormId := fmt.Sprint(oldSettings[configId].LoggingFormID)

		newConf.Forms[logFormId] = NewFormFields{
			LoggingField509A:           fmt.Sprint(oldSettings[configId].LoggingField509A),
			LoggingFieldLookupDate:     fmt.Sprint(oldSettings[configId].LoggingFieldLookupDate),
			LoggingFieldPub78:          fmt.Sprint(oldSettings[configId].LoggingFieldPub78),
			LoggingFieldRulingDate:     fmt.Sprint(oldSettings[configId].LoggingFieldRulingDate),
			LoggingFieldSubsectionDesc: fmt.Sprint(oldSettings[configId].LoggingFieldSubsectionDesc),
		}
	}
	return newConf, nil
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

	var oldConfig OldConfig
	oldRef := client.NewRef(testWs)
	oldRef.Get(ctx, &oldConfig)
	if err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	newConfig, err := migrate(oldConfig)
	if err != nil {
		log.Fatalln("Failed to convert old plugin config to new", err)
	}
	// updatedRef := client.NewRef(testWs2)
	// var updatedConfig map[string]interface{}

	// updatedRef.Get(ctx, &updatedConfig)
	// if err != nil {
	// 	log.Fatalln("Error reading from database:", err)
	// }

	oldTest, _ := json.MarshalIndent(oldConfig, "", "  ")
	newTest, _ := json.MarshalIndent(newConfig, "", "  ")
	// updatedTest, _ := json.MarshalIndent(updatedConfig, "", "  ")
	fmt.Println("Old: ", string(oldTest))
	fmt.Println("New: ", string(newTest))
	// fmt.Println(string(updatedTest))

}
