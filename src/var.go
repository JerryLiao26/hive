package main

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
var supportedCommands = [...]Command{HELP, START, GEN, SET, LIST, AUTH, ADD, DEL}

// Supported command Handler
var supportedCommandHandlers = [...]func(){helpHandler, startHandler, genHandler, setHandler, listHandler, authCliHandler, addHandler, delHandler}

// LogLevel is to defined logging level
type LogLevel string

// INFO level
const INFO LogLevel = "INFO"

// WARN level
const WARN LogLevel = "WARN"

// ERROR level
const ERROR LogLevel = "ERROR"

// Config path
var confPath struct {
	homePath string
	confFile string
	confDir  string
}

// Database config
var dbConf struct {
	username string
	password string
	addr     string
	port     string
}

// Serve config
var serveConf struct {
	addr string
	port string
}

// Cli config
var cliConf struct {
	admin string
}

// serverInfo stores server info
var serverInfo struct {
	startTime  string
	serverOS   string
	serverArch string
}

// Config file
type confItem string

// dbUsername defines database username in config file
const dbUsername confItem = "DB_USERNAME"

// dbPassword defines database password in config file
const dbPassword confItem = "DB_PASSWORD"

// dbAddress defines database address in config file
const dbAddress confItem = "DB_ADDRESS"

// dbPort defines database port in config file
const dbPort confItem = "DB_PORT"

// serveAddress defines serve address in config File
const serveAddress confItem = "SERVE_ADDRESS"

// serveAddress defines serve address in config File
const servePort confItem = "SERVE_PORT"

// authAdmin defines current authorized admin on cli
const authAdmin confItem = "AUTH_ADMIN"

// Config file items
var confItems = [...]confItem{dbUsername, dbPassword, dbAddress, dbPort, serveAddress, servePort, authAdmin}

// Config file matcher
var confItemMatcher = [...]*string{&dbConf.username, &dbConf.password, &dbConf.addr, &dbConf.port, &serveConf.addr, &serveConf.port, &cliConf.admin}
