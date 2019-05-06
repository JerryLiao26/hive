package cli

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/JerryLiao26/hive/config"
	"github.com/JerryLiao26/hive/database"
	"github.com/JerryLiao26/hive/helper"
	"github.com/JerryLiao26/hive/server"
)

// Supported command Handler
var SupportedCommandHandlers = [...]func(){HelpHandler, startHandler, genHandler, setHandler, listHandler, authCliHandler, addHandler, delHandler}

func serveHandler() {
	// TODO: Server daemon
}

func startHandler() {
	if len(os.Args) >= 3 {
		str := os.Args[2]
		group := strings.Split(str, ":")
		helper.ServeConf.Addr = group[0]
		helper.ServeConf.Port = group[1]
		config.SaveConf()
	}

	server.Serve()
}

func genHandler() {
	if helper.CliConf.Admin != "" {
		if len(os.Args) >= 3 {
			// Get data
			tag := os.Args[2]
			plain := tag + "hive"
			// Check duplicate
			if !database.CheckTagDuplicate(tag) {
				// Crypt with MD5
				obj := md5.New()
				obj.Write([]byte(plain))
				cipher := obj.Sum(nil)
				token := hex.EncodeToString(cipher)
				// Store generated token
				if database.StoreToken(tag, token) {
					helper.CliLogger("Generated token:" + token)
					helper.CliLogger("Token with tag \"" + tag + "\" stored successfully")
				} else {
					helper.CliLogger("Token store failed")
				}
			} else {
				helper.CliLogger("Tag duplicated. Please run \"hive list\" to see stored tags")
			}
		} else {
			HelpHandler()
		}
	} else {
		helper.CliLogger("No authorized admin, please use hive auth [token] to authorize")
	}
}

func setHandler() {
	if len(os.Args) >= 3 {
		str := os.Args[2]
		group := strings.Split(str, ":")
		username := group[0]
		password := group[1]
		// Set conf
		helper.DbConf.Username = username
		helper.DbConf.Password = password
		helper.DbConf.Addr = "127.0.0.1"
		helper.DbConf.Port = "3306"
		// Store conf
		config.SaveConf()
		helper.CliLogger("Database config stored")
	} else {
		HelpHandler()
	}
}

func HelpHandler() {
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
	if helper.CliConf.Admin != "" {
		output := database.FetchToken()
		helper.CliLogger("Stored token are listed below as \"tag:token\":")
		for i := 0; i < len(output); i++ {
			fmt.Println(output[i])
		}
	} else {
		helper.CliLogger("No authorized admin, please use hive auth [token] to authorize")
	}
}

func authCliHandler() {
	if len(os.Args) >= 3 {
		token := os.Args[2]
		name, flag := database.FetchAdmin(token)
		if flag {
			helper.CliConf.Admin = name
			config.SaveConf() // Save authorized admin
			helper.CliLogger("Admin \"" + name + "\" authorized successfully")
		} else {
			helper.CliLogger("No such admin with token \"" + token + "\"")
		}
	} else {
		HelpHandler()
	}
}

func addHandler() {
	if len(os.Args) >= 3 {
		// Get data
		admin := os.Args[2]
		plain := admin + "hive"
		// Check duplicate
		if !database.CheckAdminDuplicate(admin) {
			// Crypt with MD5
			obj := md5.New()
			obj.Write([]byte(plain))
			cipher := obj.Sum(nil)
			token := hex.EncodeToString(cipher)
			// Store generated token
			if database.StoreAdmin(admin, token) {
				helper.CliLogger("Generated token:" + token)
				helper.CliLogger("Admin \"" + admin + "\" stored successfully")
			} else {
				helper.CliLogger("Token store failed")
			}
		} else {
			helper.CliLogger("Admin duplicated.")
		}
	} else {
		HelpHandler()
	}
}

func delHandler() {
	if helper.CliConf.Admin != "" {
		if len(os.Args) >= 3 {
			tag := os.Args[2]
			if database.CheckTagDuplicate(tag) {
				if database.DelToken(tag) {
					helper.CliLogger("Token with tag \"" + tag + "\" deleted")
				} else {
					helper.CliLogger("Token delete failed")
				}
			} else {
				helper.CliLogger("No such tag named \"" + tag + "\". Please run \"hive list\" to see stored tags")
			}
		} else {
			HelpHandler()
		}
	} else {
		helper.CliLogger("No authorized admin, please use hive auth [token] to authorize")
	}
}

func FirstHandler() {
	// Prompt
	fmt.Println("It's your first time using hive-cli, let's configure it first.")
	// Database config
	fmt.Print("Database username:")
	_, _ = fmt.Scanln(&helper.DbConf.Username)
	fmt.Print("Database password:")
	_, _ = fmt.Scanln(&helper.DbConf.Password)
	helper.DbConf.Addr = "127.0.0.1"
	helper.DbConf.Port = "3306"
	// Server config
	fmt.Print("Serving address(without port):")
	_, _ = fmt.Scanln(&helper.ServeConf.Addr)
	fmt.Print("Serving port:")
	_, _ = fmt.Scanln(&helper.ServeConf.Port)
	// Build directory
	dirString := helper.ConfPath.HomePath + string(os.PathSeparator) + helper.ConfPath.ConfDir
	_ = os.MkdirAll(dirString, 0755)
	config.SaveConf()
	// Finish
	fmt.Println("You can now use \"hive add [name]\" to add an admin and enjoy using hive!")
}
