package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	address := host + ":" + port

	router := http.NewServeMux()
	router.HandleFunc("/", handleIndex)
	router.HandleFunc("/subscribe", handleSubscribe)

	log.Println("Starting server on", address)
	http.ListenAndServe(address, router)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `
		<html>
		<head>
		  <title>Main page</title>
		  <script src="https://unpkg.com/htmx.org@1.7.0"></script>
		  <link rel="stylesheet" href="https://the.missing.style">
		</head>
		<body>
		<header>
			<h1>Hello</h1>
		</header>
		<main>
			A simple landing page.
			<figure>
				<figcaption>Subscribe to our newsletter</figcaption>
				<form hx-get="/subscribe">
					<label>Email address <input name="email" type="email"></Label>
					<button type="submit">Subscribe</button>
				</form>
			</figure>
		</main>
		<div style="position: fixed; right: 2em; bottom: 2em; max-height: 50%; overflow: hidden;" >
			<a hx-get="/chat" hx-swap="outerHTML">Live chat</a>
		</div>
		</body>
		</html>
	`)
}

func handleSubscribe(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	fmt.Fprintf(w, `
		Confirmation email sent to %v
	`, email)
}
