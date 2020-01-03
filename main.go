package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func weebhook(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("secret") != os.Getenv("FORUNE_COOKIE") {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Page not found"))
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var a action
	err = json.Unmarshal(data, &a)
	if err != nil {
		panic(err)
	}
	switch a.QueryResult.Action {
	case "input.welcome":
		log.Println("input.welcome triggered")
		rand.Seed(time.Now().UTC().UnixNano())
		if rand.Float32() > 0.5 {
			err := sendMessage(w, NewWelcomeMessage())
			if err != nil {
				panic(err)
			}
		} else {
			err := sendMessage(w, NewEvent("custom_welcome_event"))
			if err != nil {
				panic(err)
			}
		}

	case "input.feeling":
		log.Println("input.feeling triggered")
		err := sendQuoteWithFeeling(w, a.QueryResult.Parameters.Feeling)
		if err != nil {
			panic(err)
		}

	case "input.fortune":
		log.Println("input.fortune triggered")
		err := sendQuote(w)
		if err != nil {
			panic(err)
		}

	case "input.authors":
		log.Println("input.authors triggered")
		err := sendMessage(w, NewAuthorList())
		if err != nil {
			panic(err)
		}

	case "input.author.quote":
		log.Println("input.authors.quote triggered")
		a, err := getAuthor(a.QueryResult.OutputContexts)
		if err != nil {
			panic(err)
		}
		err = sendAuthorQuote(w, a)
		if err != nil {
			panic(err)
		}

	default:
		log.Println("no input action triggered")
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/webhook", weebhook).Methods("POST")

	log.Println("Start listening on :9090...")
	err := http.ListenAndServe(":9090", r)
	if err != nil {
		panic(err)
	}
}
