package helper

// Command is to defined supported command
type Command string

// SERVE command
const SERVE Command = "serve"

// START command
const START Command = "start"

// GEN command
const GEN Command = "gen"

// SET command
const SET Command = "set"

// HELP command
const HELP Command = "help"

// LIST command
const LIST Command = "list"

// AUTH command
const AUTH Command = "auth"

// ADD command
const ADD Command = "add"

// DEL command
const DEL Command = "del"

// Supported commands
var SupportedCommands = [...]Command{HELP, START, GEN, SET, LIST, AUTH, ADD, DEL}

// LogLevel is to defined logging level
type LogLevel string

// INFO level
const INFO LogLevel = "INFO"

// WARN level
const WARN LogLevel = "WARN"

// ERROR level
const ERROR LogLevel = "ERROR"

// Config path
var ConfPath struct {
	HomePath string
	ConfFile string
	ConfDir  string
}

// Database config
var DbConf struct {
	Username string
	Password string
	Addr     string
	Port     string
}

// Serve config
var ServeConf struct {
	Addr string
	Port string
}

// Cli config
var CliConf struct {
	Admin string
}

// serverInfo stores server info
var ServerInfo struct {
	StartTime  string
	ServerOS   string
	ServerArch string
}

// message defines struct to store message data from database
type Message struct {
	ID        int    `json:"id"`
	Tag       string `json:"tag"`
	Admin     string `json:"admin"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// Config file
type ConfItem string

// dbUsername defines database username in config file
const dbUsername ConfItem = "DB_USERNAME"

// dbPassword defines database password in config file
const dbPassword ConfItem = "DB_PASSWORD"

// dbAddress defines database address in config file
const dbAddress ConfItem = "DB_ADDRESS"

// dbPort defines database port in config file
const dbPort ConfItem = "DB_PORT"

// serveAddress defines serve address in config File
const serveAddress ConfItem = "SERVE_ADDRESS"

// serveAddress defines serve address in config File
const servePort ConfItem = "SERVE_PORT"

// authAdmin defines current authorized admin on cli
const authAdmin ConfItem = "AUTH_ADMIN"

// Config file items
var ConfItems = [...]ConfItem{dbUsername, dbPassword, dbAddress, dbPort, serveAddress, servePort, authAdmin}

// Config file matcher
var ConfItemMatcher = [...]*string{&DbConf.Username, &DbConf.Password, &DbConf.Addr, &DbConf.Port, &ServeConf.Addr, &ServeConf.Port, &CliConf.Admin}
