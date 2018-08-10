package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"time"
)

// GET method
const GET string = "GET"

// POST method
const POST string = "POST"

// tokenAuth is for token authorize
type tokenAuth struct {
	Token string `json:"token"`
}

// sessionValidate defines session validate request
type sessionValidate struct {
	SessionID string `json:"sessionId"`
}

// message defines struct to store message data from database
type message struct {
	ID        int    `json:"id"`
	Tag       string `json:"tag"`
	Admin     string `json:"admin"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// messageRequest defines message request struct
type messageRequest struct {
	Text  string `json:"text"`
	Token string `json:"token"`
}

// messageRespond defines message respond struct
type messageRespond struct {
	Code     int       `json:"code"`
	Text     string    `json:"text"`
	Method   string    `json:"method"`
	Messages []message `json:"messages"`
}

// infoRespond defines respond struct for request
type infoRespond struct {
	Code       int    `json:"code"`
	Text       string `json:"text"`
	Method     string `json:"method"`
	StartTime  string `json:"startTime"`
	ServerOS   string `json:"serverOS"`
	ServerArch string `json:"serverArch"`
	AdminName  string `json:"adminName"`
}

// authRespond defines respond struct after auth
type authRespond struct {
	Code      int    `json:"code"`
	Text      string `json:"text"`
	Method    string `json:"method"`
	Name      string `json:"name"`
	SessionID string `json:"sessionId"`
}

// serverRespond defines respond struct
type serverRespond struct {
	Code   int    `json:"code"`
	Text   string `json:"text"`
	Method string `json:"method"`
}

func enterLog(w http.ResponseWriter, r *http.Request) {
	serverLogger("From", r.RemoteAddr, INFO)
	serverLogger(r.Method, r.URL.Path, INFO)
	serverLogger("Scheme", r.URL.Scheme, INFO)
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
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
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
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
			sessionID := addSession(name, ta.Token, r)
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
			serverLogger("Auth token invalid", ta.Token, WARN)
		}
	} else {
		methodNotAllowed(w, r)
		serverLogger("Auth warning", r.Method+" method is not allowed for /auth, abandoned", WARN)
	}
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == POST {
		var sv sessionValidate
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&sv)
		if err != nil {
			serverLogger("JSON parse error", err.Error(), ERROR)
		}
		name, flag := validateSession(sv.SessionID, r)
		if flag {
			group := fetchMessages(name)
			// Respond to client
			mr := messageRespond{
				Code:     200,
				Text:     "Session verified",
				Method:   r.Method,
				Messages: group,
			}
			output, err := json.Marshal(mr)
			if err != nil {
				serverLogger("JSON build error", err.Error(), ERROR)
			}
			fmt.Fprintf(w, string(output))
			serverLogger("Messages delivered", r.RemoteAddr, INFO)
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
			serverLogger("Session invalid", sv.SessionID, WARN)
		}
	} else {
		methodNotAllowed(w, r)
		serverLogger("Message warning", r.Method+" method is not allowed for /messages, abandoned", WARN)
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
		name, flag := validateSession(sv.SessionID, r)
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
			serverLogger("Session of admin \""+name+"\" verified", sv.SessionID, INFO)
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
			serverLogger("Session invalid", sv.SessionID, WARN)
		}
	} else {
		methodNotAllowed(w, r)
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
		methodNotAllowed(w, r)
		serverLogger("Panel warning", r.Method+" method is not allowed for /panel, abandoned", WARN)
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == POST {
		var sv sessionValidate
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&sv)
		if err != nil {
			serverLogger("JSON parse error", err.Error(), ERROR)
		}
		name, flag := validateSession(sv.SessionID, r)
		if flag {
			ir := infoRespond{
				Code:       200,
				Text:       "Session verified",
				Method:     r.Method,
				StartTime:  serverInfo.startTime,
				ServerOS:   serverInfo.serverOS,
				ServerArch: serverInfo.serverArch,
				AdminName:  name,
			}
			output, err := json.Marshal(ir)
			if err != nil {
				serverLogger("Session invalid", sv.SessionID, ERROR)
			}
			fmt.Fprintf(w, string(output))
			serverLogger("Info sent", r.RemoteAddr, INFO)
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
			serverLogger("Session invalid", sv.SessionID, WARN)
		}
	} else {

		serverLogger("Info warning", r.Method+" method is not allowed for /info, abandoned", WARN)
	}
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == POST {
		var m messageRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&m)
		if err != nil {
			serverLogger("JSON parse error", err.Error(), ERROR)
		}
		// Checking token
		name, tag := checkToken(m.Token)
		if tag == "" {
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
			// Storing data
			if storeMessage(name, tag, m.Text) {
				//  Response to client
				sr := serverRespond{
					Code:   200,
					Method: r.Method,
					Text:   "Message from " + tag + " saved",
				}
				output, err := json.Marshal(sr)
				if err != nil {
					serverLogger("JSON build error", err.Error(), ERROR)
				}
				fmt.Fprintf(w, string(output))
				serverLogger("Message from", tag, INFO)
			} else {
				//  Response to client
				sr := serverRespond{
					Code:   500,
					Method: r.Method,
					Text:   "Message from " + tag + " could not save",
				}
				output, err := json.Marshal(sr)
				if err != nil {
					serverLogger("JSON build error", err.Error(), ERROR)
				}
				fmt.Fprintf(w, string(output))
				serverLogger("Message saving failed", tag, WARN)
			}
		}
	} else {
		methodNotAllowed(w, r)
		serverLogger("Send warning", r.Method+" method is not allowed for /send, abandoned", WARN)
	}
}

func serve() {
	// Get serve string
	serveString := serveConf.addr + ":" + serveConf.port
	// Handlers
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/panel", panelHandler)
	http.HandleFunc("/session", sessionHandler)
	http.HandleFunc("/messages", messagesHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	// Static files
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/", fileServer)
	// Server start
	serverInfo.startTime = time.Now().Format("2006-01-02 15:04:05")
	serverInfo.serverOS = runtime.GOOS
	serverInfo.serverArch = runtime.GOARCH
	serverLogger("Starting", "Serve at "+serveString, INFO)
	err := http.ListenAndServe(serveString, nil)
	// Error
	if err != nil {
		serverLogger("Cannot start", err.Error(), ERROR)
	}
}
