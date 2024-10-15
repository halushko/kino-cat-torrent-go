package listeners

import (
	"context"
	"fmt"
	"github.com/hekmon/transmissionrpc/v2"
	"kino-cat-torrent-go/helpers"
	"log"
	"math"
	"strings"
)

func GetAllTorrents() {
	processor := func(key string, args []string, client *transmissionrpc.Client) string {
		log.Printf("[GetAllTorrents] Торенти отримано")
		torrents, err := client.TorrentGetAll(context.Background())
		if err != nil {
			text := fmt.Sprintf("[GetAllTorrents] Помилка отримання переліку торентов: %v", err)
			log.Printf(text)
			return text
		}
		log.Printf("[GetAllTorrents] Торенти для сзовища %s отримано", key)
		answer := generateAnswerList(args[0], args[1], torrents)
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_LIST", processor)
}

func generateAnswerList(server string, id string, torrents []transmissionrpc.Torrent) string {
	var line strings.Builder

	for _, torrent := range torrents {
		line.WriteString(fmt.Sprintf("%s %s\n", getStatusIcon(torrent), *torrent.Name))
		line.WriteString(fmt.Sprintf("%s %s\n", getProgressBar(torrent), getGigabytesLeft(torrent)))
		line.WriteString(fmt.Sprintf("/more_%s_%s ", server, id))
		line.WriteString(fmt.Sprintf("/files_%s_%s\n", server, id))
	}
	return line.String()
}

func getProgressBar(torrent transmissionrpc.Torrent) string {
	blocks := 20
	blackBlocks := int(*torrent.PercentDone * float64(blocks))
	var line strings.Builder
	line.WriteString("||")
	for i := 0; i < blackBlocks; i++ {
		line.WriteString("█")
	}
	if blackBlocks < blocks {
		line.WriteString("▒")
	}
	if blackBlocks+1 < blocks {
		for i := blackBlocks + 1; i < blocks; i++ {
			line.WriteString("░")
		}
	}
	return line.String()
}

func getGigabytesLeft(torrent transmissionrpc.Torrent) string {
	done := *torrent.PercentDone
	if done == 1.0 {
		return " (заверш)"
	}

	percentDone := math.Round(done * 100)
	totalSize := float64(*torrent.TotalSize)
	remainingSize := math.Round((totalSize-(totalSize*done))/1000000.0) / 1000.0

	return fmt.Sprintf("%.0f %% (%.3f Gb залиш)", percentDone, remainingSize)
}

func getStatusIcon(torrent transmissionrpc.Torrent) string {
	switch *torrent.Status {
	case transmissionrpc.TorrentStatusStopped:
		return "⏸"
	case transmissionrpc.TorrentStatusCheckWait:
		return "⏲♾"
	case transmissionrpc.TorrentStatusCheck:
		return "♾"
	case transmissionrpc.TorrentStatusDownloadWait:
		return "⏲⬇️"
	case transmissionrpc.TorrentStatusDownload:
		return "⬇️"
	case transmissionrpc.TorrentStatusSeedWait:
		return "⏲⬆️"
	case transmissionrpc.TorrentStatusSeed:
		return "⬆️"
	case transmissionrpc.TorrentStatusIsolated:
		return "🈲"
	default:
		return "🈲"
	}
}
