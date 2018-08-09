package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

// GET method
const GET string = "GET"

// POST method
const POST string = "POST"

// message defines message struct
type message struct {
	Text  string `json:"text"`
	Token string `json:"token"`
}

// tokenAuth is for token authorize
type tokenAuth struct {
	Token string `json:"token"`
}

// authRespond defines respond struct after auth
type authRespond struct {
	Code      int    `json:"code"`
	Text      string `json:"text"`
	Method    string `json:"method"`
	Name      string `json:"name"`
	SessionID string `json:"sessionId"`
}

// sessionValidate defines session validate request
type sessionValidate struct {
	SessionID string `json:"sessionId"`
}

// serverRespond defines respond struct
type serverRespond struct {
	Code   int    `json:"code"`
	Text   string `json:"text"`
	Method string `json:"method"`
}

// sendRespond defines respond struct
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
	fmt.Println("request", r)
	t := template.Must(template.ParseFiles("template/hello.html"))
	t.Execute(w, serveConf.port)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == GET {
		t := template.Must(template.ParseFiles("template/auth.html"))
		t.Execute(w, nil)
	} else if r.Method == POST {
		var ta tokenAuth
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&ta)
		if err != nil {
			serverLogger("JSON parse error", err.Error(), ERROR)
		}
		// Validate token
		name, flag := fetchAdmin(ta.Token)
		if flag {
			// Generate session
			sessionID := addSession(ta.Token, r)
			//  Response to client
			ar := authRespond{
				Code:      200,
				Method:    r.Method,
				Text:      "Token verified",
				Name:      name,
				SessionID: sessionID,
			}
			output, err := json.Marshal(ar)
			if err != nil {
				serverLogger("JSON build error", err.Error(), ERROR)
			}
			fmt.Fprintf(w, string(output))
			serverLogger("Auth token verified", ta.Token, INFO)
		} else {
			//  Response to client
			sr := serverRespond{
				Code:   500,
				Method: r.Method,
				Text:   "Token invalid",
			}
			output, err := json.Marshal(sr)
			if err != nil {
				serverLogger("JSON build error", err.Error(), ERROR)
			}
			fmt.Fprintf(w, string(output))
			serverLogger("Auth token invalid", ta.Token, ERROR)
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
		serverLogger("Auth warning", r.Method+" method is not allowed for /auth, abandoned", WARN)
	}
}

func sessionHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == POST {
		var sv sessionValidate
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&sv)
		if err != nil {
			serverLogger("JSON parse error", err.Error(), ERROR)
		}
		flag := validateSession(sv.SessionID, r)
		if flag {
			//  Response to client
			sr := serverRespond{
				Code:   200,
				Method: r.Method,
				Text:   "Session verified",
			}
			output, err := json.Marshal(sr)
			if err != nil {
				serverLogger("JSON build error", err.Error(), ERROR)
			}
			fmt.Fprintf(w, string(output))
			serverLogger("Session verified", sv.SessionID, INFO)
		} else {
			//  Response to client
			sr := serverRespond{
				Code:   500,
				Method: r.Method,
				Text:   "Session invalid",
			}
			output, err := json.Marshal(sr)
			if err != nil {
				serverLogger("JSON build error", err.Error(), ERROR)
			}
			fmt.Fprintf(w, string(output))
			serverLogger("Session invalid", sv.SessionID, ERROR)
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
		serverLogger("Session warning", r.Method+" method is not allowed for /session, abandoned", WARN)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == GET {
		t := template.Must(template.ParseFiles("template/dashboard.html"))
		t.Execute(w, nil)
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
		serverLogger("Dashboard warning", r.Method+" method is not allowed for /dashboard, abandoned", WARN)
	}
}

func panelHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == GET {
		t := template.Must(template.ParseFiles("template/panel.html"))
		t.Execute(w, nil)
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
		serverLogger("Panel warning", r.Method+" method is not allowed for /panel, abandoned", WARN)
	}
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == POST {
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
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/session", sessionHandler)
	http.HandleFunc("/panel", panelHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
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
