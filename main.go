package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Repository Repository `json:"repository"`
	Ref        string     `json:"ref"`
	Before     string     `json:"before"`
	Forced     bool       `json:"forced"`
}

type Repository struct {
	FullName string `json:"full_name"`
	HtmlUrl  string `json:"html_url"`
}

func main() {
	godotenv.Load(".env")
	router := mux.NewRouter().StrictSlash(true)
	router.
		Methods("POST").
		Path(os.Getenv("path_to_hook")).
		HandlerFunc(pushEvent)

	log.Fatal(http.ListenAndServe(os.Getenv("host")+":"+os.Getenv("port"), router))
}

func pushEvent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var msg Message
	json.Unmarshal(body, &msg)
	json, _ := json.Marshal(msg.Repository)
	f, err := os.OpenFile(os.Getenv("json_file"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	if _, err = f.WriteString(string(json) + "\n"); err != nil {
		log.Println(err)
	}
}
