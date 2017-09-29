package main

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"bufio"
	"fmt"
)

var filePointer int64
func initWatcher() {
	filePointer = 0
}
func fileModified() {
	stat,err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	if filePointer > stat.Size() {
		filePointer = 0
		insertPosition(0)
	}

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i, err := getDocument(scanner.Text())
		if err == nil {
			insertLog(i)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	filePointer = stat.Size()

}
func WatchFile(filename string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	done := make(chan bool)

	go func(){
		for {
			select {
				case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Modified")
					fileModified()
				}
			}
		}
	}()

	err = watcher.Add(filename)
	if err != nil {
		panic(err)
	}
	<-done
}