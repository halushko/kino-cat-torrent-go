package listeners

import (
	"context"
	"fmt"
	"github.com/hekmon/transmissionrpc/v2"
	"kino-cat-torrent-go/helpers"
	"log"
	"math"
	"sort"
	"strings"
)

func GetAllTorrents() {
	processor := func(args []string, client *transmissionrpc.Client) string {
		log.Printf("[GetAllTorrents] Торенти отримано")
		torrents, err := client.TorrentGetAll(context.Background())
		if err != nil {
			text := fmt.Sprintf("[GetAllTorrents] Помилка отримання переліку торентов: %v", err)
			log.Printf(text)
			return text
		}
		sort.Slice(torrents, func(i, j int) bool { return *torrents[i].ID < *torrents[j].ID })
		log.Printf("[GetAllTorrents] Торенти отримано")
		answer := generateAnswerList(torrents)
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_LIST", processor)
}

func generateAnswerList(torrents []transmissionrpc.Torrent) string {
	var line strings.Builder
	for _, torrent := range torrents {
		id := *torrent.ID
		line.WriteString(fmt.Sprintf("%s %s\n", getStatusIcon(torrent), *torrent.Name))
		line.WriteString(fmt.Sprintf("%s %s\n", getProgressBar(*torrent.PercentDone, 20), getGigabytesLeft(torrent)))
		line.WriteString(fmt.Sprintf("/more_%d ", id))
		line.WriteString(fmt.Sprintf("/files_%d\n", id))
	}
	return line.String()
}

func getProgressBar(percentDone float64, blocks int) string {

	blackBlocks := int(percentDone * float64(blocks))
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
	line.WriteString("||")
	return line.String()
}

func getGigabytesLeft(torrent transmissionrpc.Torrent) string {
	done := *torrent.PercentDone
	if done == 1.0 {
		return " (заверш)"
	}

	percentDone := math.Round(done * 100)
	totalSize := torrent.TotalSize.GB()
	remainingSize := totalSize * (1.0 - done)

	return fmt.Sprintf("%.0f %% (%.2f Gb залиш)", percentDone, remainingSize)
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
