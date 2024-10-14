package main

import (
	"github.com/halushko/kino-cat-core-go/logger_helper"
	"kino-cat-torrent-go/listeners"
	"log"
)

//goland:noinspection ALL
func main() {
	logFile := logger_helper.SoftPrepareLogFile()
	log.SetOutput(logFile)

	listeners.GetAllTorrents()

	defer logger_helper.SoftLogClose(logFile)
	select {}
}
