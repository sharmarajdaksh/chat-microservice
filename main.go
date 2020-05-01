package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func main() {
	const port = 3001

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal("[ERROR]: ", err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("[CONNECT]: ", s.ID())
		return nil
	})

	server.OnEvent("/", "message", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "[MESSAGE]: " + msg
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("[ERROR]: ", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("[DISCONNECT]: ", reason)
	})

	// Run as non-blocking thread
	go server.Serve()
	defer server.Close()

	// Root path '/'
	// The server handles on socket.io connections
	http.Handle("/", server)

	log.Println(fmt.Sprintf("[STARTUP] Serving at localhost:%d...", port))
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}
}
