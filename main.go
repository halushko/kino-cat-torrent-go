package main

import (
	"github.com/halushko/kino-cat-core-go/logger_helper"
	"kino-cat-torrent-go/listeners"
)

//goland:noinspection ALL
func main() {
	logger_helper.SoftPrepareLogFile()

	listeners.GetAllTorrents()
	listeners.GetMoreCommands()
	listeners.ExecutePauseTorrent()
	listeners.ExecuteResumeTorrent()
	listeners.GetTorrentInfo()

	select {}
}
