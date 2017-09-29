package main

import (
	"flag"
	"fmt"
)

const LOG_PATH  = "/var/log/nginx/access.log"
var filePath string

func init() {
	initDB()
	initParser()
	initWatcher()
}

func main() {
	pathPtr := flag.String("path", LOG_PATH, "nginx access log path")
	addUser := flag.Bool("adduser", false, "add a new user")

	flag.Parse()
	filePath = *pathPtr
	if *addUser == true {
		var username, password string
		fmt.Print("Enter New User > ")
		fmt.Scanln(&username)
		fmt.Print("Enter New Password > ")
		fmt.Scanln(&password)
		insertUser(username, password)
	}
	WatchFile(*pathPtr)
}