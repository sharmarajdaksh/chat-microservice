package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// socketReader struct
type socketReader struct {
	conn *websocket.Conn
	name string
}

var socketreader []*socketReader

func handler(w http.ResponseWriter, r *http.Request) {
	if socketreader == nil {
		socketreader = make([]*socketReader, 0)
	}

	defer func() {
		err := recover()
		if err != nil {
			log.Println("[ERR] ", err)
		}
		r.Body.Close()
	}()

	// Upgrader upgrades the http connection to a websocket connection
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// Upgrade connection
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERR] ", err)
	}

	// Socket Reader
	ptrSocketReader := &socketReader{
		conn: con,
		name: "",
	}

	socketreader = append(socketreader, ptrSocketReader)

	// Start a new thread
	ptrSocketReader.startThread()
}

func (i *socketReader) broadcast(str string) {
	for _, g := range socketreader {
		g.writeMsg(i.name, str)
	}
}

func (i *socketReader) read() {
	_, b, er := i.conn.ReadMessage()
	if er != nil {
		panic(er)
	}

	i.broadcast(string(b))

	log.Println("[MESSAGE] " + i.name + " [SAYS] " + string(b))
}

func (i *socketReader) writeMsg(name string, str string) {
	i.conn.WriteMessage(websocket.TextMessage, []byte(name+"--->"+str))
}

func (i *socketReader) startThread() {

	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Println("[ERR]", err)
			}
		}()

		for {
			if i.name == "" {
				_, b, err := i.conn.ReadMessage()
				if err != nil {
					log.Println("[ERR]", err)
				}
				i.name = string(b)
				log.Println("[JOINED]", i.name)
			} else {
				i.read()
			}
		}
	}()
}

func main() {

	h := http.NewServeMux()
	h.HandleFunc("/", handler)

	http.ListenAndServe(":3001", h)
}
