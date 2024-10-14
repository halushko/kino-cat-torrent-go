package main

import (
	"github.com/halushko/kino-cat-core-go/logger_helper"
	"log"
)

//goland:noinspection ALL
func main() {
	logFile := logger_helper.SoftPrepareLogFile()
	log.SetOutput(logFile)

	//listeners.StartUserMessageListener()
	//listeners.StartGetHelpCommandListener()

	block := make(chan struct{})
	defer logger_helper.SoftLogClose(logFile)
	<-block
}
