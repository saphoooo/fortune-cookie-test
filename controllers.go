package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewEvent creates a new event message
func NewEvent(text string) *Response {
	return &Response{
		FollowupEventInput: &followupEventInput{
			Name: text,
		},
	}
}

// NewMessage creates a new message
func NewMessage(text string) *Response {
	return &Response{
		FulfillmentText: text,
	}
}

// NewAuthorList creates an author list
func NewAuthorList() *Response {
	return &Response{
		Payload: &payload{
			Google: &google{
				ExpectUserResponse: true,
				RichResponse: &richResponse{
					Items: []items{
						items{
							SimpleResponse: &simpleResponse{
								TextToSpeech: "Select an author",
							},
						},
					},
				},
				SystemIntent: &systemIntent{
					Intent: "actions.intent.OPTION",
					Data: &data{
						Type: "type.googleapis.com/google.actions.v2.OptionValueSpec",
						ListSelect: &listSelect{
							Title: "Select an author",
							ListItems: []listItems{
								listItems{
									OptionInfo: &optionInfo{
										Key:      "Eliot",
										Synonyms: []string{"Eliot", "T. S. Eliot"},
									},
									Title: "T. S. Eliot",
								},
								listItems{
									OptionInfo: &optionInfo{
										Key:      "White",
										Synonyms: []string{"James", "White", "James White"},
									},
									Title: "James White",
								},
								listItems{
									OptionInfo: &optionInfo{
										Key:      "Stutman",
										Synonyms: []string{"Dave", "Stutman", "Dave Stutman"},
									},
									Title: "Dave Stutman",
								},
								listItems{
									OptionInfo: &optionInfo{
										Key:      "Churchill",
										Synonyms: []string{"Winston", "Churchill", "Winston Churchill"},
									},
									Title: "Winston Churchill",
								},
								listItems{
									OptionInfo: &optionInfo{
										Key:      "Allen",
										Synonyms: []string{"Woody Allen", "Woody", "Allen"},
									},
									Title: "Woody Allen",
								},
								listItems{
									OptionInfo: &optionInfo{
										Key:      "Twain",
										Synonyms: []string{"Mark", "Twain", "Mark Twain"},
									},
									Title: "Mark Twain",
								},
								listItems{
									OptionInfo: &optionInfo{
										Key:      "Einstein",
										Synonyms: []string{"Albert", "Einstein", "Albert Einstein"},
									},
									Title: "Albert Einstein",
								},
								listItems{
									OptionInfo: &optionInfo{
										Key:      "Wright",
										Synonyms: []string{"Steven", "Wright", "Steven Wright"},
									},
									Title: "Steven Wright",
								},
							},
						},
					},
				},
			},
		},
	}
}

func sendMessage(w http.ResponseWriter, msg interface{}) error {
	resp, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	log.Println(string(resp))
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
	return nil
}

func sendQuote(w http.ResponseWriter) error {
	rand.Seed(time.Now().UTC().UnixNano())
	quote, err := findQuoteByNumber(rand.Intn(8))
	if err != nil {
		return err
	}
	err = sendMessage(w, NewMessage(quote))
	if err != nil {
		return err
	}
	return nil
}

func sendQuoteWithFeeling(w http.ResponseWriter, parameter string) error {
	switch parameter {
	case "happy":
		rand.Seed(time.Now().UTC().UnixNano())
		quote, err := findQuoteByNumber(rand.Intn(5))
		if err != nil {
			return err
		}
		err = sendMessage(w, NewMessage(quote))
		if err != nil {
			return err
		}
		return nil
	case "sad":
		rand.Seed(time.Now().UTC().UnixNano())
		quote, err := findQuoteByNumber(rand.Intn(3) + 5)
		if err != nil {
			return err
		}
		err = sendMessage(w, NewMessage(quote))
		if err != nil {
			return err
		}
		return nil
	default:
		rand.Seed(time.Now().UTC().UnixNano())
		quote, err := findQuoteByNumber(rand.Intn(8))
		if err != nil {
			return err
		}
		err = sendMessage(w, NewMessage(quote))
		if err != nil {
			return err
		}
		return nil
	}
}

func sendAuthorQuote(w http.ResponseWriter, author string) error {
	q, err := findQuoteByAuthor(author)
	if err != nil {
		return err
	}
	if q == "" {
		return errors.New("no match found")
	}
	err = sendMessage(w, NewMessage(q))
	if err != nil {
		return err
	}
	return nil
}

func getAuthor(o []outputContexts) (string, error) {
	for _, v := range o {
		if v.Parameters.Option != "" {
			return v.Parameters.Option, nil
		}
	}
	return "", errors.New("author not found")
}

// NewWelcomeMessage ...
func NewWelcomeMessage() *Response {
	return &Response{
		Payload: &payload{
			Google: &google{
				RichResponse: &richResponse{
					Items: []items{
						items{
							SimpleResponse: &simpleResponse{
								Ssml: `<speak>
								<prosody rate="high" pitch="+4st">
								"Hello, welcome to Zguingou's fortune cookie"</prosody>
								</speak>`,
								DisplayText: "Hello, welcome to Zguingou's fortune cookie",
							},
						},
					},
				},
			},
		},
	}
}

func findQuoteByAuthor(author string) (string, error) {
	var result struct {
		Quote string
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collection := client.Database("fortuneCookie").Collection("quotes")
	filter := bson.M{"synonym": author}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.Quote, nil
}

func findQuoteByNumber(number int) (string, error) {
	var a []quotes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collection := client.Database("fortuneCookie").Collection("quotes")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return "", err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result quotes
		err := cur.Decode(&result)
		if err != nil {
			return "", err
		}
		a = append(a, result)
	}
	if err := cur.Err(); err != nil {
		return "", err
	}
	return a[number].Quote, nil
}
