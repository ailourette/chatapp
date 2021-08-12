package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/nsqio/go-nsq"
)

type Message struct {
	Name      string
	Content   string
	Timestamp string
}

var (
	tpl *template.Template
)

func sendMsg(msgName, msgContent string) error {
	//The only valid way to instantiate the Config
	config := nsq.NewConfig()
	//Creating the Producer using NSQD Address
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}
	//Init topic name and message
	topic := "Topic_Example"
	msg := Message{
		Name:      msgName,
		Content:   msgContent,
		Timestamp: time.Now().String(),
	}
	//Convert message as []byte
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	//Publish the Message
	err = producer.Publish(topic, payload)
	if err != nil {
		//log.Println(err)
		return err
	}
	return nil
}

func getMsg(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		validate := validator.New()
		msgName := req.FormValue("msgName")
		err1 := validate.Var(msgName, "required")
		msgContent := req.FormValue("msgContent")
		err2 := validate.Var(msgContent, "required")

		if err1 != nil || err2 != nil {
			http.Error(res, "Please enter your message name and message content!", http.StatusForbidden)
			return
		} else {
			err := sendMsg(msgName, msgContent)
			if err == nil {
				io.WriteString(res, `
				<html>
				<meta http-equiv='refresh' content='5; url=/sendMsg '/>
				Message sent successfully!<br>
				You will be redirected shortly in 5 seconds...<br>
				</html>
			`)
				return
			} else {
				http.Error(res, "Message send unsuccessfully!", http.StatusForbidden)
				return
			}
		}
	}
	tpl.ExecuteTemplate(res, "getMsg.gohtml", nil)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/sendMsg", getMsg)
}
