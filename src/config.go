package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func initPath() string {
	// Get pathString
	confPath.homePath = os.Getenv("HOME")
	confPath.confDir = ".hive"
	confPath.confFile = "hive.conf"
	return confPath.homePath + string(os.PathSeparator) + confPath.confDir + string(os.PathSeparator) + confPath.confFile
}

func loadConf() {
	pathString := initPath()
	inFile, err := os.Open(pathString)
	if err != nil {
		// File not exist
		if os.IsNotExist(err) {
			// First-time using
			firstHandler()
			return
		}
		cliLogger(err.Error())
	}
	// Close file
	defer inFile.Close()
	// Load config
	reader := bufio.NewReader(inFile)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			// Check if file end
			if err == io.EOF {
				return
			}
			cliLogger(err.Error())
		}
		// Restore config
		confLine := strings.Split(str, ":")
		for i := 0; i < len(confItems); i++ {
			if confItems[i] == confItem(confLine[0]) {
				length := len(confLine[1])
				*confItemMatcher[i] = confLine[1][:length-1]
			}
		}
	}
}

func saveConf() {
	pathString := initPath()
	outFile, err := os.OpenFile(pathString, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		cliLogger(err.Error())
	}
	// Close file
	defer outFile.Close()
	// Save config
	writer := bufio.NewWriter(outFile)
	// Store config
	for i := 0; i < len(confItems); i++ {
		str := string(confItems[i]) + ":" + *confItemMatcher[i] + "\n"
		writer.WriteString(str)
	}
	writer.Flush()
}
