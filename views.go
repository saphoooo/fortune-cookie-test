package main

type quotes struct {
	Author  string `json:"author,omitempty"`
	Synonym string `json:"synonym,omitempty"`
	Tags    string `json:"tags,omitempty"`
	Quote   string `json:"quote,omitempty"`
}

type action struct {
	ResponseID  string `json:"responseId,omitempty"`
	QueryResult struct {
		QueryText  string `json:"queryText,omitempty"`
		Action     string `json:"action,omitempty"`
		Parameters struct {
			Date    string `json:"date,omitempty"`
			Feeling string `json:"Feeling,omitempty"`
		} `json:"parameters,omitempty"`
		AllRequiredParamsPresent bool `json:"allRequiredParamsPresent,omitempty"`
		FulfillmentMessages      []struct {
			Text struct {
				Text []string `json:"text,omitempty"`
			} `json:"text,omitempty"`
		} `json:"fulfillmentMessages,omitempty"`
		OutputContexts []outputContexts `json:"outputContexts,omitempty"`
		Intent         struct {
			Name        string `json:"name,omitempty"`
			DisplayName string `json:"displayName,omitempty"`
		} `json:"intent,omitempty"`
		IntentDetectionConfidence float32 `json:"intentDetectionConfidence,omitempty"`
		LanguageCode              string  `json:"languageCode,omitempty"`
	} `json:"queryResult,omitempty"`
	OriginalDetectIntentRequest struct {
		Payload struct{} `json:"payload,omitempty"`
	} `json:"originalDetectIntentRequest,omitempty"`
	Session string `json:"session,omitempty"`
}

type outputContexts struct {
	Parameters parameters `json:"parameters,omitempty"`
}

type parameters struct {
	Option string `json:"OPTION,omitempty"`
}

// Response ...
type Response struct {
	FulfillmentText    string              `json:"fulfillmentText,omitempty"`
	FollowupEventInput *followupEventInput `json:"followupEventInput,omitempty"`
	Payload            *payload            `json:"payload,omitempty"`
}

type followupEventInput struct {
	Name string `json:"name,omitempty"`
}

type payload struct {
	Google *google `json:"google,omitempty"`
}

type google struct {
	ExpectUserResponse bool          `json:"expectUserResponse,omitempty"`
	RichResponse       *richResponse `json:"richResponse,omitempty"`
	SystemIntent       *systemIntent `json:"systemIntent,omitempty"`
}

type richResponse struct {
	Items []items `json:"items,omitempty"`
}

type items struct {
	SimpleResponse *simpleResponse `json:"simpleResponse,omitempty"`
}

type simpleResponse struct {
	TextToSpeech string `json:"textToSpeech,omitempty"`
	Ssml         string `json:"ssml,omitempty"`
	DisplayText  string `json:"displayText,omitempty"`
}

type systemIntent struct {
	Intent string `json:"intent,omitempty"`
	Data   *data  `json:"data,omitempty"`
}

type data struct {
	Type       string      `json:"@type,omitempty"`
	ListSelect *listSelect `json:"listSelect,omitempty"`
}

type listSelect struct {
	Title     string      `json:"title,omitempty"`
	ListItems []listItems `json:"items,omitempty"`
}

type listItems struct {
	OptionInfo *optionInfo `json:"optionInfo,omitempty"`
	Title      string      `json:"title,omitempty"`
}

type optionInfo struct {
	Key      string   `json:"key,omitempty"`
	Synonyms []string `json:"synonyms,omitempty"`
}
