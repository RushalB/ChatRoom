package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

type Authentication struct {
	Username string
	Status   bool
}

type User struct {
	Username string
	Password string
}

type Client struct {
	Username string
	Conn     *websocket.Conn
}

type Chat struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}
type Message1 struct {

	// 0 ==> Chat message
	// 1 ==> User Data

	Type int    `json:"type"`
	Body string `json:"body"`
}

type Message struct {
	// 2 ==> Errors
	// -1 ==> Previous load
	// 0 ==> False -- Recvd
	// 1 ==> True -- Sent

	Type int    `json:"type"`
	Body []Chat `json:"body"`
}

var chat_messages = []Chat{}
var clients = []Client{}
var count = 0

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func notifyAll(c Chat) {
	for _, cl := range clients {
		message := Message{Type: 2, Body: []Chat{(c)}}
		if err := cl.Conn.WriteJSON(message); err != nil {
			fmt.Println("6")
			log.Println(err)
			return
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"Hello")
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		//Process the request
		var u User
		reqBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(reqBody, &u)
		// fmt.Println(u.Username)

		//Send Response
		var auth Authentication
		auth.Username = u.Username

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		if u.Password == "b" {
			auth.Status = true
		} else {
			auth.Status = false

		}
		json.NewEncoder(w).Encode(auth)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Start(cli Client) {

	var conn = cli.Conn

	for {
		var m Message1
		err := conn.ReadJSON(&m)
		if err != nil {
			fmt.Println(err)
			for i := 0; i < len(clients); i++ {
				if clients[i].Conn == conn {
					index_to_remove := i
					clients = append(clients[:index_to_remove], clients[index_to_remove+1:]...)
				}
			}

			log.Println(cli.Username + " left the seission")

			chat := Chat{Sender: cli.Username, Text: string(cli.Username + " left the seission")}
			notifyAll(chat)
			return
		}
		if m.Type == 1 {
			cli.Username = m.Body
			c := Chat{Sender: cli.Username, Text: string(cli.Username + " joined the seission")}
			notifyAll(c)
			log.Println(cli.Username + " joined the seission")
			// return
		} else {

			//broadcast here-loop through client list and send WriteMessage()[easy]
			//for every read from client there is a write to all

			chat := Chat{Sender: cli.Username, Text: string(m.Body)}
			chat_messages = append(chat_messages, chat)
			for _, cl := range clients {
				message := Message{Type: 1, Body: []Chat{chat}}
				if cl.Conn == conn {
					message.Type = 1
				} else {
					message.Type = 0
				}
				if err := cl.Conn.WriteJSON(message); err != nil {
					fmt.Println("3")
					log.Println(err)
					return
				}
			}
		}
	}
}

func WS(w http.ResponseWriter, r *http.Request) {

	//Hello Client
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("0")
		log.Println(err)
	}
	//Register your client

	client := Client{
		Username: "User" + fmt.Sprint(count),
		Conn:     ws}
	count += 2
	//Add client to List of clients
	clients = append(clients, client)

	//Load client with message history
	if len(chat_messages) > 0 {

		message := Message{Type: -1, Body: chat_messages}
		err1 := ws.WriteJSON(message)
		if err1 != nil {
			fmt.Println("1")
			log.Println(err1)
		}
	}

	//Start Client
	Start(client)

}

func setup() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/login", login)
	http.HandleFunc("/ws", WS)
	http.ListenAndServe(":8000", nil)
}

func main() {
	setup()

}
