package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
)

/**
 * This is the main package, the starting point for the application
 */

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxActive:   80,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

var (
	Pool        *redis.Pool
	redisServer = flag.String("localhost", ":4245", "")
)

func init() {
	flag.Parse()
	Pool = newPool(*redisServer)
}

/**
 * Handler for static files
 */
func FileServerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := vars["filename"]
	fType := vars["type"]
	filepath := "/home/jordan/convo/static/" + fType + "/" + file
	//the file path here may need some refactoring on different environments
	http.ServeFile(w, r, filepath)

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Card struct {
	Type string `json:"type"`
	Text string `json:"text"`
	User string `json:"user"`
	Date int    `json:"date"`
}

/*HomeHandler : handler for home*/
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/jordan/convo/static/html/home.html")
}

/*SocketHandler :
 * Handles socket connections
 */
func SocketHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go RedisPubSubHandler(conn)

	v := Card{}

	for {
		err := conn.ReadJSON(&v)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(v)
		if v.Type == "test" {
			testMessage := Card{
				Type: "test",
				Text: "Connection Successful...",
				User: "server",
			}
			if err = conn.WriteJSON(testMessage); err != nil {
				log.Println(err)
				return
			}
		} else {
			PublishMessage(v)
		}
	}
}

func RedisPubSubHandler(socketConn *websocket.Conn) {
	var newMessage Card
	var newMessageString string
	var newMessageBytes []byte

	c := Pool.Get()
	defer c.Close()

	psc := redis.PubSubConn{c}
	psc.Subscribe("main")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
			newMessageString = string(v.Data)
			newMessageBytes = []byte(newMessageString)
			err := json.Unmarshal(newMessageBytes, &newMessage)
			if err != nil {
				log.Println(err)
				break
			}
			if err = socketConn.WriteJSON(newMessage); err != nil {
				log.Println(err)
				return
			}
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			return
		}
	}
}

func PublishMessage(message Card) {
	var messageString string
	var messageBytes []byte
	var err error

	c := Pool.Get()
	defer c.Close()

	messageBytes, err = json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	messageString = string(messageBytes)

	c.Do("PUBLISH", "main", messageString)
}

func main() {

	log.Println("Server running on port 8080")

	r := mux.NewRouter()

	// Main Handlers
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/static/{type}/{filename}", FileServerHandler)
	r.HandleFunc("/socket", SocketHandler)

	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
