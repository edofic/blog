package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := host + ":" + port

	chat := NewChat()

	router := http.NewServeMux()
	router.HandleFunc("/", handleIndex)
	router.HandleFunc("/chat", handleChat)
	router.HandleFunc("/chat/events", newHandleEvents(chat))
	router.HandleFunc("/chat/msg", newHandleMsg(chat))

	log.Println("Starting server on", address)
	http.ListenAndServe(address, router)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `
		<html>
		<head>
		  <title>Chat</title>
		  <script src="https://unpkg.com/htmx.org@1.7.0"></script>
		  <link rel="stylesheet" href="https://the.missing.style">
		</head>
		<body>
		  <main hx-trigger="load" hx-get="/chat" />
		</body>
		</html>
		`)
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `
		<div class="chat-main">
			<style>
			.chat-main {
			  background: lightyellow;
			  padding: 1em;
			}
			</style>
			<h3>Chat</h3>
			<div hx-trigger="load" hx-get="/chat/msg"></div>
			<div hx-sse="connect:/chat/events swap:message" hx-swap="afterbegin"></div>
		</div>
	`)
}

func newHandleEvents(chat *Chat) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		msgs, unsubscribe := chat.Register()
		defer unsubscribe()

		for msg := range msgs {
			fmt.Fprintf(w, "data: <div>%v</div>\n\n", msg)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

const msgTemplate = `
<form hx-post="/chat/msg">
<input type="text" name="msg" autocomplete="off">
<button type="submit">Send</button>
<form>
`

func newHandleMsg(chat *Chat) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := r.FormValue("msg")
		if msg != "" {
			chat.Broadcast(msg)
		}
		fmt.Fprintln(w, msgTemplate)
	}
}

type Chat struct {
	mu      sync.Mutex
	clients map[int]chan<- string
	counter int
}

func NewChat() *Chat {
	return &Chat{
		clients: make(map[int]chan<- string),
	}
}

func (c *Chat) Register() (<-chan string, func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch := make(chan string, 1)
	key := c.counter
	c.counter += 1
	c.clients[key] = ch
	return ch, func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.clients, key)
	}
}

func (c *Chat) Broadcast(msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var toDelete []int
	for key, ch := range c.clients {
		select {
		case ch <- msg:
		default:
			close(ch)
			toDelete = append(toDelete, key)
		}
	}
	for _, key := range toDelete {
		delete(c.clients, key)
	}
}
