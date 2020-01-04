package main

import (
	"context"
	"fmt"
	"os"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// DetectIntentText ...
func DetectIntentText(projectID, sessionID, text, languageCode string) (string, error) {
	ctx := context.Background()

	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return "", err
	}
	defer sessionClient.Close()

	if projectID == "" || sessionID == "" {
		return "", fmt.Errorf("Received empty project (%s) or session (%s)", projectID, sessionID)
	}

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: languageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return "", err
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()
	return fulfillmentText, nil
}

// DoThisNotThat ...
func DoThisNotThat(projectID, sessionID, text string) ([]string, error) {
	languageCode := "en-US"
	ctx := context.Background()

	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return nil, err
	}
	defer sessionClient.Close()

	if projectID == "" || sessionID == "" {
		return nil, fmt.Errorf("Received empty project (%s) or session (%s)", projectID, sessionID)
	}

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: languageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return nil, err
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()
	queryIntent := queryResult.GetIntent()
	queryText := queryResult.GetQueryText()
	var a []string
	a = append(a, fulfillmentText, queryIntent.DisplayName, queryText)
	return a, nil
}

func main() {
	projectID := os.Getenv("FORTUNE_COOKIE_PROJECTID")
	sessionID := "my-test-session-id"
	//languageCode := "en-US"
	text := "Give me a fortune cookie"
	m, err := DoThisNotThat(projectID, sessionID, text)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Detected intent\n\tQuery: %s\n\tResponse: %s\n\tIntent: %s\n", m[2], m[0], m[1])
}
