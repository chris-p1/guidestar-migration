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
	"os"

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

type NewConfigFields struct {
	Enabled                    bool   `json:"enabled"`
	LoggingField509A           string `json:"loggingField509a,omitempty"`
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
	Forms    map[string]FormFields  `json:"forms,omitempty"`
	Settings map[string]interface{} `json:"settings"`
}

type Firebase map[string]Config

func convert(oldFields map[string]OldConfigFields) Config {

	newConf := Config{}
	newFormFields := map[string]FormFields{}
	newSettings := make(map[string]interface{})

	for k, v := range oldFields {

		formId := fmt.Sprint(oldFields[k].LoggingFormID)
		newFormFields[formId] = FormFields{
			LoggingField509A:           fmt.Sprint(v.LoggingField509A),
			LoggingFieldLookupDate:     fmt.Sprint(v.LoggingFieldLookupDate),
			LoggingFieldPub78:          fmt.Sprint(v.LoggingFieldPub78),
			LoggingFieldRulingDate:     fmt.Sprint(v.LoggingFieldRulingDate),
			LoggingFieldSubsectionDesc: fmt.Sprint(v.LoggingFieldSubsectionDesc),
		}

		newSettings[k] = NewConfigFields{
			Enabled:                    true,
			LoggingField509A:           fmt.Sprint(v.LoggingField509A),
			LoggingFieldLookupDate:     fmt.Sprint(v.LoggingFieldLookupDate),
			LoggingFieldPub78:          fmt.Sprint(v.LoggingFieldPub78),
			LoggingFieldRulingDate:     fmt.Sprint(v.LoggingFieldRulingDate),
			LoggingFieldSubsectionDesc: fmt.Sprint(v.LoggingFieldSubsectionDesc),
			LoggingFormEnabled:         true,
			LoggingLinkedFormID:        fmt.Sprint(v.LoggingFieldEin),
			LoggingFormID:              v.LoggingFormID,
			Mch1:                       Mch1{Type: "Logging", Value: "Yes"},
			Name:                       v.Name,
			TargetFieldID:              v.TargetFieldID,
			TargetFormID:               v.TargetFormID,
		}

	}
	newConf.Forms = newFormFields
	newConf.Settings = newSettings

	// fmt.Println(newConf)
	return newConf
}

func main() {
	fmt.Println(reverse.String("Guidestar Migration"))

	fbUrl := "https://guidestar-multiconfig.firebaseio.com"
	opt := option.WithCredentialsFile("config/guidestar-multiconfig-firebase-adminsdk-pzpvm-0fcc9de3b9.json")
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

	newConfig := Firebase{}
	var fbConfig Firebase
	ref := client.NewRef("")
	ref.Get(ctx, &fbConfig)
	if err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	for k := range fbConfig {

		// only update old plugin settings
		if len(fbConfig[k].Forms) == 0 {

			var old map[string]OldConfigFields
			buf, _ := json.Marshal(fbConfig[k].Settings)
			if err := json.Unmarshal(buf, &old); err != nil {
				log.Fatal(err)
			}
			wsConf := convert(old)
			newConfig[k] = wsConf

		} else {

			newConfig[k] = fbConfig[k]
		}
	}
	oldTest, _ := json.MarshalIndent(fbConfig, "", "  ")
	newTest, _ := json.MarshalIndent(newConfig, "", "  ")

	fmt.Println("Old: ", string(oldTest))
	fmt.Println("New: ", string(newTest))

	err = os.WriteFile("oldTest.json", oldTest, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("newTest.json", newTest, 0644)
	if err != nil {
		log.Fatal(err)
	}

}
