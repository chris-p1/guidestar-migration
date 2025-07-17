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

type FormFields struct {
	LoggingField509A           string `json:"loggingField509a"`
	LoggingFieldLookupDate     string `json:"loggingFieldLookupDate"`
	LoggingFieldPub78          string `json:"loggingFieldPub78"`
	LoggingFieldRulingDate     string `json:"loggingFieldRulingDate"`
	LoggingFieldSubsectionDesc string `json:"loggingFieldSubsectionDesc"`
}

type ConfigFields struct {
	Enabled                    bool   `json:"enabled,omitempty"`
	LoggingField509A           string `json:"loggingField509a,omitempty"`
	LoggingFieldEin            int    `json:"loggingFieldEin,omitempty"`
	LoggingFieldLookupDate     string `json:"loggingFieldLookupDate,omitempty"`
	LoggingFieldPub78          string `json:"loggingFieldPub78,omitempty"`
	LoggingFieldRulingDate     string `json:"loggingFieldRulingDate,omitempty"`
	LoggingFieldSubsectionDesc string `json:"loggingFieldSubsectionDesc,omitempty"`
	LoggingFormEnabled         bool   `json:"loggingFormEnabled,omitempty"`
	LoggingLinkedFormID        string `json:"loggingLinkedFormId,omitempty"`
	LoggingFormID              int    `json:"loggingFormId,omitempty"`
	Mch1                       Mch1   `json:"mch1,omitzero"`
	Name                       string `json:"name"`
	TargetFieldID              int    `json:"targetFieldId"`
	TargetFormID               int    `json:"targetFormId"`
}

type Mch1 struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type Config struct {
	Forms    map[string]FormFields   `json:"forms,omitempty"`
	Settings map[string]ConfigFields `json:"settings,omitempty"`
}

type Firebase map[string]Config

func migrate(oldConfig Config) (Config, error) {
	// var newConfig map[string]interface{}
	newConf := Config{
		Forms:    make(map[string]FormFields),
		Settings: make(map[string]ConfigFields),
	}
	oldSettings := oldConfig.Settings

	for configId := range oldSettings {

		fmt.Printf("Key: %s\nValue: %#v\n", configId, oldSettings[configId])

		newConf.Settings[configId] = ConfigFields{
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

		newConf.Forms[logFormId] = FormFields{
			LoggingField509A:           oldSettings[configId].LoggingField509A,
			LoggingFieldLookupDate:     oldSettings[configId].LoggingFieldLookupDate,
			LoggingFieldPub78:          oldSettings[configId].LoggingFieldPub78,
			LoggingFieldRulingDate:     oldSettings[configId].LoggingFieldRulingDate,
			LoggingFieldSubsectionDesc: oldSettings[configId].LoggingFieldSubsectionDesc,
		}
	}

	return newConf, nil
}

// func checkNewConf(newSettings *NewConfigFields, newForms *NewFormFields) {

// }

func main() {
	fmt.Println(reverse.String("Guidestar Migration"))

	// fileContent, err := os.Open("guidestar-multiconfig-export.json")
	// if err != nil {
	// 	log.Fatalf("Error reading JSON file: %v", err)
	// }

	// defer fileContent.Close()
	// byteValue, _ := io.ReadAll(fileContent)
	// var oldFb Firebase
	// err = json.Unmarshal(byteValue, &oldFb)
	// if err != nil {
	// 	log.Fatalf("Error unmarshaling")
	// }

	// // fmt.Println(oldFb)
	// // var newFb Firebase
	// for wsId := range oldFb {
	// 	fmt.Printf("Workspace: %s\n", wsId)
	// 	fmt.Println(oldFb[wsId])
	// }
	// 20250714T114955 TODO: change to prod once ready
	// testWs := "33940"
	// // testWs2 := "40138"
	// // testWsProd := "38099"
	fbUrl := "https://guidestar-multiconfig.firebaseio.com"
	opt := option.WithCredentialsFile("config/guidestar-multiconfig-firebase-adminsdk-pzpvm-0fcc9de3b9.json")
	// // fbUrl := "https://guidestar-stage.firebaseio.com"
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

	var oldConfig Firebase
	oldRef := client.NewRef("")
	oldRef.Get(ctx, &oldConfig)
	if err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	// newConfig, err := migrate(oldConfig)
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
	// newTest, _ := json.MarshalIndent(newConfig, "", "  ")
	// updatedTest, _ := json.MarshalIndent(updatedConfig, "", "  ")
	fmt.Println("Old: ", string(oldTest))
	// fmt.Println("New: ", string(newTest))
	// fmt.Println(string(updatedTest))

}
