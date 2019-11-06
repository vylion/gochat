package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Message struct
type Text struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

type Message struct {
	Id      string `json:"id"`
	Content Text   `json:"content"`
}

// Chat as a slice of Messages
var chat []Message

// Possible Messages
var quickchat []string

//Possible authors
var users []string

var idCount int

func getChat(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(chat)
}
func getMsg(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for _, msg := range chat {
		if msg.Id == params["id"] {
			json.NewEncoder(writer).Encode(msg)
			return
		}
	}
}

func createMsg() (string, string) {
	content := quickchat[rand.Intn(len(quickchat))]
	author := users[rand.Intn(len(users))]

	return content, author
}

func sendMsg(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var newMsg Message
	_ = json.NewDecoder(request.Body).Decode(&(newMsg.Content))
	newMsg.Id = strconv.Itoa(idCount)
	idCount++
	chat = append(chat, newMsg)
	json.NewEncoder(writer).Encode(newMsg)
	fmt.Println("New Message:")
	fmt.Println("#", newMsg.Id, newMsg.Content.Author, ":", newMsg.Content.Content)
}
func editMsg(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	params := mux.Vars(request)
	for _, msg := range chat {
		if msg.Id == params["id"] {
			_ = json.NewDecoder(request.Body).Decode(&(msg.Content))
			msg.Content.Content = "(Edited) " + msg.Content.Content
			fmt.Println("Edited Message:")
			fmt.Println("#", msg.Id, msg.Content.Author, ":", msg.Content.Content)
			//chat = append(chat[:index], chat[index+1:]...)
			json.NewEncoder(writer).Encode(chat)
			return
		}
	}
}
func deleteMsg(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	params := mux.Vars(request)
	for index, msg := range chat {
		if msg.Id == params["id"] {
			chat = append(chat[:index], chat[index+1:]...)
			break
		}
	}
	json.NewEncoder(writer).Encode(chat)
}

func main() {
	// Init rand
	rand.Seed(time.Now().Unix())

	// Init server router
	router := mux.NewRouter()

	// Init chat variables
	quickchat = []string{"Wow!", "OMG!", "Nice shot!", "What a save!", "Sorry!", "My bad."}
	users = []string{"Chipper", "Beast", "Jester", "Merlin", "Bandit", "Rainmaker"}

	first := Message{"0", Text{"The chat room is open.", "SERVER"}}
	chat = append(chat, first)
	idCount = 1

	// Handlers
	router.HandleFunc("/api/chat", getChat).Methods("GET")
	router.HandleFunc("/api/chat/{id}", getMsg).Methods("GET")
	router.HandleFunc("/api/chat", sendMsg).Methods("POST")
	router.HandleFunc("/api/chat/{id}", editMsg).Methods("PUT")
	router.HandleFunc("/api/chat/{id}", deleteMsg).Methods("DELETE")

	// Log
	log.Fatal(http.ListenAndServe(":8000", router))

	last := Message{"-1", Text{"The chat room is closed.", "SERVER"}}
	chat = append(chat, last)
}
