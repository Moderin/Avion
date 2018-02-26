package main

import (
	"bufio"
	"color"
	"fmt"
	"os"
	"os/user"
	"strings"
)

var config = make(map[string]*string)

func configFileSave() {
	var configFile *os.File
	configFile, err := os.Create("config")
	if err != nil {
		fmt.Println(color.Red("Cannot make config file "), err)
	}
	defer configFile.Close()

	for key, value := range config {
		line := key + ":" + *value
		writer := bufio.NewWriter(configFile)
		fmt.Fprintln(writer, line)
		writer.Flush()
	}
}

func loadConfig() {
	configFile, err := os.Open("config")
	if err == nil {
		scanner := bufio.NewScanner(configFile)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			slices := strings.Split(scanner.Text(), ":")
			config[slices[0]] = &slices[1]
		}
	}

	// set default values if empty

	if config["name"] == nil {
		currentUser, _ := user.Current()
		config["name"] = &currentUser.Name
	}

	if config["avatar"] == nil {
		path := "img/avatar.png"
		config["avatar"] = &path
	}
	configFile.Close()
}
