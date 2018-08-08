package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type message struct {
	Text  string `json:"text"`
	Token string `json:"token"`
}

type serverRespond struct {
	Code   int    `json:"code"`
	Text   string `json:"text"`
	Method string `json:"method"`
}

type sendRespond struct {
	Code   int    `json:"code"`
	Tag    string `json:"tag"`
	Text   string `json:"text"`
	Sender string `json:"sender"`
}

func enterLog(w http.ResponseWriter, r *http.Request) {
	serverLogger("From", r.RemoteAddr, INFO)
	serverLogger(r.Method, r.URL.Path, INFO)
	serverLogger("Scheme", r.URL.Scheme, INFO)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	t := template.Must(template.ParseFiles("template/hello.html"))
	t.Execute(w, serveConf.port)
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == "POST" {
		var m message
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&m)
		if err != nil {
			serverLogger("JSON parse error", err.Error(), ERROR)
		}
		// Checking token
		res := checkToken(m.Token)
		if res == "" {
			//  Response to client
			sr := serverRespond{
				Code:   400,
				Method: r.Method,
				Text:   "Token invalid",
			}
			output, err := json.Marshal(sr)
			if err != nil {
				serverLogger("JSON build error", err.Error(), ERROR)
			}
			fmt.Fprintf(w, string(output))
			serverLogger("Token invalid", m.Token, WARN)
		} else {
			user := res
			// Storing data
			if storeMessage(user, m.Text) {
				//  Response to client
				sr := serverRespond{
					Code:   200,
					Method: r.Method,
					Text:   "Message from " + user + " saved",
				}
				output, err := json.Marshal(sr)
				if err != nil {
					serverLogger("JSON build error", err.Error(), ERROR)
				}
				fmt.Fprintf(w, string(output))
				serverLogger("Message from", user, INFO)
			} else {
				//  Response to client
				sr := serverRespond{
					Code:   500,
					Method: r.Method,
					Text:   "Message from " + user + " could not save",
				}
				output, err := json.Marshal(sr)
				if err != nil {
					serverLogger("JSON build error", err.Error(), ERROR)
				}
				fmt.Fprintf(w, string(output))
				serverLogger("Message saving failed", user, ERROR)
			}
		}
	} else {
		sr := serverRespond{
			Code:   400,
			Method: r.Method,
			Text:   "Method not allowed",
		}
		output, err := json.Marshal(sr)
		if err != nil {
			serverLogger("JSON build error", err.Error(), ERROR)
		}
		fmt.Fprintf(w, string(output))
		serverLogger("Send warning", r.Method+" method is not allowed for /send, abandoned", WARN)
	}
}

func serve() {
	// Get serve string
	serveString := serveConf.addr + ":" + serveConf.port
	// Handlers
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/hello", helloHandler)
	// Static files
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/", fileServer)
	// Server start
	serverLogger("Starting", "Serve at "+serveString, INFO)
	err := http.ListenAndServe(serveString, nil)
	// Error
	if err != nil {
		serverLogger("Cannot start", err.Error(), ERROR)
	}
}
