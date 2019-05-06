package config

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/JerryLiao26/hive/helper"
)

func initPath() string {
	// Get pathString
	helper.ConfPath.HomePath = os.Getenv("HOME")
	helper.ConfPath.ConfDir = ".hive"
	helper.ConfPath.ConfFile = "hive.conf"
	return helper.ConfPath.HomePath + string(os.PathSeparator) + helper.ConfPath.ConfDir + string(os.PathSeparator) + helper.ConfPath.ConfFile
}

func LoadConf() bool {
	pathString := initPath()
	inFile, err := os.Open(pathString)
	if err != nil {
		// File not exist
		if os.IsNotExist(err) {
			// First-time using
			return false
		}
		helper.CliLogger(err.Error())
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
				return true
			}
			helper.CliLogger(err.Error())
		}
		// Restore config
		confLine := strings.Split(str, ":")
		for i := 0; i < len(helper.ConfItems); i++ {
			if helper.ConfItems[i] == helper.ConfItem(confLine[0]) {
				length := len(confLine[1])
				*helper.ConfItemMatcher[i] = confLine[1][:length-1]
			}
		}
	}
}

func SaveConf() {
	pathString := initPath()
	outFile, err := os.OpenFile(pathString, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		helper.CliLogger(err.Error())
	}
	// Close file
	defer outFile.Close()
	// Save config
	writer := bufio.NewWriter(outFile)
	// Store config
	for i := 0; i < len(helper.ConfItems); i++ {
		str := string(helper.ConfItems[i]) + ":" + *helper.ConfItemMatcher[i] + "\n"
		_, _ = writer.WriteString(str)
	}
	_ = writer.Flush()
}
