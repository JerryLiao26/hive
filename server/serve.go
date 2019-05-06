package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"time"

	"github.com/JerryLiao26/hive/database"
	"github.com/JerryLiao26/hive/helper"
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

// messageRequest defines message request struct
type messageRequest struct {
	Text  string `json:"text"`
	Token string `json:"token"`
}

// messageRespond defines message respond struct
type messageRespond struct {
	Code     int              `json:"code"`
	Text     string           `json:"text"`
	Method   string           `json:"method"`
	Messages []helper.Message `json:"messages"`
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
	helper.ServerLogger("From", r.RemoteAddr, helper.INFO)
	helper.ServerLogger(r.Method, r.URL.Path, helper.INFO)
	helper.ServerLogger("Scheme", r.URL.Scheme, helper.INFO)
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	sr := serverRespond{
		Code:   400,
		Method: r.Method,
		Text:   "Method not allowed",
	}
	output, err := json.Marshal(sr)
	if err != nil {
		helper.ServerLogger("JSON build error", err.Error(), helper.ERROR)
	}
	_, _ = fmt.Fprintf(w, string(output))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	t := template.Must(template.ParseFiles("template/hello.html"))
	_ = t.Execute(w, helper.ServeConf.Port)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == GET {
		t := template.Must(template.ParseFiles("template/auth.html"))
		_ = t.Execute(w, nil)
	} else if r.Method == POST {
		var ta tokenAuth
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&ta)
		if err != nil {
			helper.ServerLogger("JSON parse error", err.Error(), helper.ERROR)
		}
		// Validate token
		name, flag := database.FetchAdmin(ta.Token)
		if flag {
			// Generate session
			sessionID := helper.AddSession(name, ta.Token, r)
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
				helper.ServerLogger("JSON build error", err.Error(), helper.ERROR)
			}
			_, _ = fmt.Fprintf(w, string(output))
			helper.ServerLogger("Auth token verified", ta.Token, helper.INFO)
		} else {
			//  Response to client
			sr := serverRespond{
				Code:   500,
				Method: r.Method,
				Text:   "Token invalid",
			}
			output, err := json.Marshal(sr)
			if err != nil {
				helper.ServerLogger("JSON build error", err.Error(), helper.ERROR)
			}
			_, _ = fmt.Fprintf(w, string(output))
			helper.ServerLogger("Auth token invalid", ta.Token, helper.WARN)
		}
	} else {
		methodNotAllowed(w, r)
		helper.ServerLogger("Auth warning", r.Method+" method is not allowed for /auth, abandoned", helper.WARN)
	}
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == POST {
		var sv sessionValidate
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&sv)
		if err != nil {
			helper.ServerLogger("JSON parse error", err.Error(), helper.ERROR)
		}
		name, flag := helper.ValidateSession(sv.SessionID, r)
		if flag {
			group := database.FetchMessages(name)
			// Respond to client
			mr := messageRespond{
				Code:     200,
				Text:     "Session verified",
				Method:   r.Method,
				Messages: group,
			}
			output, err := json.Marshal(mr)
			if err != nil {
				helper.ServerLogger("JSON build error", err.Error(), helper.ERROR)
			}
			_, _ = fmt.Fprintf(w, string(output))
			helper.ServerLogger("Messages delivered", r.RemoteAddr, helper.INFO)
		} else {
			//  Response to client
			sr := serverRespond{
				Code:   500,
				Method: r.Method,
				Text:   "Session invalid",
			}
			output, err := json.Marshal(sr)
			if err != nil {
				helper.ServerLogger("JSON build error", err.Error(), helper.ERROR)
			}
			_, _ = fmt.Fprintf(w, string(output))
			helper.ServerLogger("Session invalid", sv.SessionID, helper.WARN)
		}
	} else {
		methodNotAllowed(w, r)
		helper.ServerLogger("Message warning", r.Method+" method is not allowed for /messages, abandoned", helper.WARN)
	}
}

func sessionHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == POST {
		var sv sessionValidate
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&sv)
		if err != nil {
			helper.ServerLogger("JSON parse error", err.Error(), helper.ERROR)
		}
		name, flag := helper.ValidateSession(sv.SessionID, r)
		if flag {
			//  Response to client
			sr := serverRespond{
				Code:   200,
				Method: r.Method,
				Text:   "Session verified",
			}
			output, err := json.Marshal(sr)
			if err != nil {
				helper.ServerLogger("JSON build error", err.Error(), helper.ERROR)
			}
			_, _ = fmt.Fprintf(w, string(output))
			helper.ServerLogger("Session of admin \""+name+"\" verified", sv.SessionID, helper.INFO)
		} else {
			//  Response to client
			sr := serverRespond{
				Code:   500,
				Method: r.Method,
				Text:   "Session invalid",
			}
			output, err := json.Marshal(sr)
			if err != nil {
				helper.ServerLogger("JSON build error", err.Error(), helper.ERROR)
			}
			_, _ = fmt.Fprintf(w, string(output))
			helper.ServerLogger("Session invalid", sv.SessionID, helper.WARN)
		}
	} else {
		methodNotAllowed(w, r)
		helper.ServerLogger("Session warning", r.Method+" method is not allowed for /session, abandoned", helper.WARN)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == GET {
		t := template.Must(template.ParseFiles("template/dashboard.html"))
		_ = t.Execute(w, nil)
	} else {
		sr := serverRespond{
			Code:   400,
			Method: r.Method,
			Text:   "Method not allowed",
		}
		output, err := json.Marshal(sr)
		if err != nil {
			helper.ServerLogger("JSON build error", err.Error(), helper.ERROR)
		}
		_, _ = fmt.Fprintf(w, string(output))
		helper.ServerLogger("Dashboard warning", r.Method+" method is not allowed for /dashboard, abandoned", helper.WARN)
	}
}

func panelHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == GET {
		t := template.Must(template.ParseFiles("template/panel.html"))
		_ = t.Execute(w, nil)
	} else {
		methodNotAllowed(w, r)
		helper.ServerLogger("Panel warning", r.Method+" method is not allowed for /panel, abandoned", helper.WARN)
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == POST {
		var sv sessionValidate
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&sv)
		if err != nil {
			helper.ServerLogger("JSON parse error", err.Error(), helper.ERROR)
		}
		name, flag := helper.ValidateSession(sv.SessionID, r)
		if flag {
			ir := infoRespond{
				Code:       200,
				Text:       "Session verified",
				Method:     r.Method,
				StartTime:  helper.ServerInfo.StartTime,
				ServerOS:   helper.ServerInfo.ServerOS,
				ServerArch: helper.ServerInfo.ServerArch,
				AdminName:  name,
			}
			output, err := json.Marshal(ir)
			if err != nil {
				helper.ServerLogger("Session invalid", sv.SessionID, helper.ERROR)
			}
			_, _ = fmt.Fprintf(w, string(output))
			helper.ServerLogger("Info sent", r.RemoteAddr, helper.INFO)
		} else {
			//  Response to client
			sr := serverRespond{
				Code:   500,
				Method: r.Method,
				Text:   "Session invalid",
			}
			output, err := json.Marshal(sr)
			if err != nil {
				helper.ServerLogger("JSON build error", err.Error(), helper.ERROR)
			}
			_, _ = fmt.Fprintf(w, string(output))
			helper.ServerLogger("Session invalid", sv.SessionID, helper.WARN)
		}
	} else {
		methodNotAllowed(w, r)
		helper.ServerLogger("Info warning", r.Method+" method is not allowed for /info, abandoned", helper.WARN)
	}
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	enterLog(w, r)
	if r.Method == POST {
		var m messageRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&m)
		if err != nil {
			helper.ServerLogger("JSON parse error", err.Error(), helper.ERROR)
		}
		// Checking token
		name, tag := database.CheckToken(m.Token)
		if tag == "" {
			//  Response to client
			sr := serverRespond{
				Code:   400,
				Method: r.Method,
				Text:   "Token invalid",
			}
			output, err := json.Marshal(sr)
			if err != nil {
				helper.ServerLogger("JSON build error", err.Error(), helper.ERROR)
			}
			_, _ = fmt.Fprintf(w, string(output))
			helper.ServerLogger("Token invalid", m.Token, helper.WARN)
		} else {
			// Storing data
			if database.StoreMessage(name, tag, m.Text) {
				//  Response to client
				sr := serverRespond{
					Code:   200,
					Method: r.Method,
					Text:   "Message from " + tag + " saved",
				}
				output, err := json.Marshal(sr)
				if err != nil {
					helper.ServerLogger("JSON build error", err.Error(), helper.ERROR)
				}
				_, _ = fmt.Fprintf(w, string(output))
				helper.ServerLogger("Message from", tag, helper.INFO)
			} else {
				//  Response to client
				sr := serverRespond{
					Code:   500,
					Method: r.Method,
					Text:   "Message from " + tag + " could not save",
				}
				output, err := json.Marshal(sr)
				if err != nil {
					helper.ServerLogger("JSON build error", err.Error(), helper.ERROR)
				}
				_, _ = fmt.Fprintf(w, string(output))
				helper.ServerLogger("Message saving failed", tag, helper.WARN)
			}
		}
	} else {
		methodNotAllowed(w, r)
		helper.ServerLogger("Send warning", r.Method+" method is not allowed for /send, abandoned", helper.WARN)
	}
}

func Serve() {
	// Get serve string
	serveString := helper.ServeConf.Addr + ":" + helper.ServeConf.Port
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
	helper.ServerInfo.StartTime = time.Now().Format("2006-01-02 15:04:05")
	helper.ServerInfo.ServerOS = runtime.GOOS
	helper.ServerInfo.ServerArch = runtime.GOARCH
	helper.ServerLogger("Starting", "Serve at "+serveString, helper.INFO)
	err := http.ListenAndServe(serveString, nil)
	// Error
	if err != nil {
		helper.ServerLogger("Cannot start", err.Error(), helper.ERROR)
	}
}
