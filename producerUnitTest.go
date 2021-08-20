package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

var (
	mainInput      int
	messageName    string
	messageContent string
)

type Message struct {
	Name      string
	Content   string
	Timestamp string
	Language  string
}

//go:generate mockgen -source=producer.go -destination=./mock/mock.go -package=mock
type publisher interface {
	Publish(topic string, body []byte) error
}

type Handler struct {
	producer publisher
}

func New(producer publisher) *Handler {
	return &Handler{
		producer: producer,
	}
}

const topic = "Topic_Example"

func (h *Handler) SendMsg(name, content string) {
	//Init topic name and message
	msg := Message{
		Name:      name,
		Content:   content,
		Timestamp: time.Now().String(),
		Language:  "es",
	}
	//Convert message as []byte
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	//Publish the Message
	err = h.producer.Publish(topic, payload)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	//The only valid way to instantiate the Config
	config := nsq.NewConfig()
	//Creating the Producer using NSQD Address
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}
	h := New(producer)
	h.SendMsg("name example", "ejemplo de contenido")
	fmt.Println("Sent")
}
