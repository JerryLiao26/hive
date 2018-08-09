package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func str2comm(str string) Command {
	return Command(str)
}

func checkSupport(comm Command) int {
	for i := 0; i < len(supportedCommands); i++ {
		if supportedCommands[i] == comm {
			return i
		}
	}
	return 0
}

func serveHandler() {
	// TODO: Server daemon
}

func startHandler() {
	if len(os.Args) >= 3 {
		str := os.Args[2]
		group := strings.Split(str, ":")
		serveConf.addr = group[0]
		serveConf.port = group[1]
		saveConf()
	}
	serve()
}

func genHandler() {
	if cliConf.admin != "" {
		if len(os.Args) >= 3 {
			// Get data
			tag := os.Args[2]
			plain := tag + "hive"
			// Check duplicate
			if !checkTagDuplicate(tag) {
				// Crypt with MD5
				obj := md5.New()
				obj.Write([]byte(plain))
				cipher := obj.Sum(nil)
				token := hex.EncodeToString(cipher)
				// Store generated token
				if storeToken(tag, token) {
					cliLogger("Generated token:" + token)
					cliLogger("Token with tag \"" + tag + "\" stored successfully")
				} else {
					cliLogger("Token store failed")
				}
			} else {
				cliLogger("Tag duplicated. Please run \"hive list\" to see stored tags")
			}
		} else {
			helpHandler()
		}
	} else {
		cliLogger("No authorized admin, please use hive auth [token] to authorize")
	}
}

func setHandler() {
	if len(os.Args) >= 3 {
		str := os.Args[2]
		group := strings.Split(str, ":")
		username := group[0]
		password := group[1]
		// Set conf
		dbConf.username = username
		dbConf.password = password
		dbConf.addr = "127.0.0.1"
		dbConf.port = "3306"
		// Store conf
		saveConf()
		cliLogger("Database config stored")
	} else {
		helpHandler()
	}
}

func helpHandler() {
	fmt.Println("Usage:")
	fmt.Println("hive [command] [args...]")
	fmt.Println("Commands:")
	// Print help text
	fmt.Println("help                       Print help text")
	fmt.Println("add [name]                 Add an admin with provided name")
	fmt.Println("auth [token]               Auth an admin to use functions")
	fmt.Println("list                       List all valid tags and corresponding tokens")
	fmt.Println("start                      Start server on last used address and port")
	// fmt.Println("serve                      Start server deamon on last used address and port")
	fmt.Println("start [address]:[port]     Start server on provided address and port")
	// fmt.Println("serve [address]:[port]     Start server deamon on provided address and port")
	fmt.Println("gen [tag]                  Generate token for provided tag")
	fmt.Println("del [tag]                  Delete token for provided tag")
	fmt.Println("set [username]:[password]  Set database connection info")
}

func listHandler() {
	if cliConf.admin != "" {
		output := fetchToken()
		cliLogger("Stored token are listed below as \"tag:token\":")
		for i := 0; i < len(output); i++ {
			fmt.Println(output[i])
		}
	} else {
		cliLogger("No authorized admin, please use hive auth [token] to authorize")
	}
}

func authCliHandler() {
	if len(os.Args) >= 3 {
		token := os.Args[2]
		name, flag := fetchAdmin(token)
		if flag {
			cliConf.admin = name
			saveConf() // Save authorized admin
			cliLogger("Admin \"" + name + "\" authorized successfully")
		} else {
			cliLogger("No such admin with token \"" + token + "\"")
		}
	} else {
		helpHandler()
	}
}

func addHandler() {
	if len(os.Args) >= 3 {
		// Get data
		admin := os.Args[2]
		plain := admin + "hive"
		// Check duplicate
		if !checkAdminDuplicate(admin) {
			// Crypt with MD5
			obj := md5.New()
			obj.Write([]byte(plain))
			cipher := obj.Sum(nil)
			token := hex.EncodeToString(cipher)
			// Store generated token
			if storeAdmin(admin, token) {
				cliLogger("Generated token:" + token)
				cliLogger("Admin \"" + admin + "\" stored successfully")
			} else {
				cliLogger("Token store failed")
			}
		} else {
			cliLogger("Admin duplicated.")
		}
	} else {
		helpHandler()
	}
}

func delHandler() {
	if cliConf.admin != "" {
		if len(os.Args) >= 3 {
			tag := os.Args[2]
			if checkTagDuplicate(tag) {
				if delToken(tag) {
					cliLogger("Token with tag \"" + tag + "\" deleted")
				} else {
					cliLogger("Token delete failed")
				}
			} else {
				cliLogger("No such tag named \"" + tag + "\". Please run \"hive list\" to see stored tags")
			}
		} else {
			helpHandler()
		}
	} else {
		cliLogger("No authorized admin, please use hive auth [token] to authorize")
	}
}

func firstHandler() {
	// Prompt
	fmt.Println("It's your first time using hive-cli, let's configure it first.")
	// Database config
	fmt.Print("Database username:")
	fmt.Scanln(&dbConf.username)
	fmt.Print("Database password:")
	fmt.Scanln(&dbConf.password)
	dbConf.addr = "127.0.0.1"
	dbConf.port = "3306"
	// Server config
	fmt.Print("Serving address(without port):")
	fmt.Scanln(&serveConf.addr)
	fmt.Print("Serving port:")
	fmt.Scanln(&serveConf.port)
	// Build directory
	dirString := confPath.homePath + string(os.PathSeparator) + confPath.confDir
	os.MkdirAll(dirString, 0755)
	saveConf()
	// Finish
	fmt.Println("You can now use \"hive add [name]\" to add an admin and enjoy using hive!")
}

func main() {
	if len(os.Args) >= 2 {
		// Pre-load
		loadConf()
		// Get command
		comm := str2comm(os.Args[1])
		// Handler for command
		supportedCommandHandlers[checkSupport(comm)]()
	} else {
		helpHandler()
	}
}
