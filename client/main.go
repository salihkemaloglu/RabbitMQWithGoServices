package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"goji.io"
	"goji.io/pat"

	"github.com/streadway/amqp"
)

type Item struct {
	Value string `bson:"value" json:"value"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!Service is working")
}

// POST a new item
func SendItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	mess, err := LoadConfiguration(item.Value)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, mess+"!:"+err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, item)
}
func LoadConfiguration(value string) (string, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
	if err != nil {
		return "Failed to connect to RabbitMQ", err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return "Failed to open a channel", err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return "Failed to declare a queue", err
	}
	body := value
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	if err != nil {
		return "Failed to publish a message", err
	}
	return "", nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func handleRequests() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/"), homePage)
	mux.HandleFunc(pat.Post("/item"), SendItem)
	log.Fatal(http.ListenAndServe(":8080", cors.AllowAll().Handler(mux)))
}

func main() {
	handleRequests()
}
