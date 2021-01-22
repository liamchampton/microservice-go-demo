package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func main() {
	// Create the route handler listening on '/'
	http.HandleFunc("/", Home)
	http.HandleFunc("/form", Form)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

// Home function is the main route for the application
func Home(w http.ResponseWriter, r *http.Request) {
	// Assign the 'msg' variable with a string value
	msg := "<h1>Client</h1>\n<p>Please use /form to send a message to the sever"

	// Write the response to the byte array - Sprintf formats and returns a string without printing it anywhere
	w.Write([]byte(fmt.Sprintf(msg)))
}

//Form parses the user input and makes POST request
func Form(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("user.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		var url string
		appRoute := os.Getenv("APPROUTE")

		fmt.Println("route", appRoute)

		// submit form logic

		// get form input
		fmt.Println("username:", r.Form["username"])

		// convert []string -> string data type
		username := strings.Join(r.Form["username"], " ")

		if appRoute != "" {
			url = appRoute + "/message"
		} else {
			url = "http://localhost:3000/message"
		}

		fmt.Println("URL:>", url)

		// POST request to the server
		json := []byte(`{"name":"` + username + `"}`)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))

		msg := "sent '" + username + "' to server"
		w.Write([]byte(fmt.Sprintf(msg)))
	}
}
