package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func init() {
	// Attempt to read the secret key as an enviroment variable
	secret := os.Getenv("SECRET_KEY_1")
	if secret == "" {
		log.Println("SECRET_KEY_1 is not set")
		return
	}
	// Create a directory for certs
	err := os.MkdirAll("./cert", 0700)
	if err != nil {
		log.Println("Error creating directory:", err)
		return
	}
	// Write this key to a file
	err = os.WriteFile("./cert/secretKey1.pem", []byte(secret), 0644)
	if err != nil {
		log.Println("Error writing file:", err)
	}
}

func main() {
	// Check for port number argument
	portPtr := flag.Int("port", 8080, "Listening Port Number")
	flag.Parse()

	// Set up the router
	r := mux.NewRouter()
	r.HandleFunc("/jearly/hello", func(w http.ResponseWriter, r *http.Request) {
		// Anonymous request
		fmt.Fprintf(w, "Hello there. Who is this?\n")
	})
	r.HandleFunc("/jearly/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		// Get the path variable
		vars := mux.Vars(r)
		name := vars["name"]
		fmt.Fprintf(w, "Hello %s!\n", name)
	})
	r.HandleFunc("/jearly/cert/secret", func(w http.ResponseWriter, r *http.Request) {
		// Show the secret file
		content, err := ioutil.ReadFile("./cert/secretKey1.pem")
		if err != nil {
			fmt.Fprintf(w, "Error reading file: %s", err.Error())
		} else {
			fmt.Fprintf(w, "%s", content)
		}
	})

	// Construct port string
	portStr := ":" + strconv.Itoa(*portPtr)
	log.Println("Listening on port", portStr)

	err := http.ListenAndServe(portStr, r)
	if err != nil {
		log.Println("Server error:", err)
	}
}
