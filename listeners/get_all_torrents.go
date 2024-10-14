package listeners

import (
	"context"
	"fmt"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"github.com/hekmon/transmissionrpc/v2"
	"github.com/nats-io/nats.go"
	"kino-cat-torrent-go/helpers"
	"log"
	"math"
	"strings"
)

func GetAllTorrents() {
	processor := func(msg *nats.Msg) {
		client, inputMessage := helpers.ConnectToTransmission(msg)
		log.Printf("[GetAllTorrents] Старт отримання переліку торентов")

		torrents, err := client.TorrentGetAll(context.Background())
		if err != nil {
			log.Printf("[GetAllTorrents] Помилка отримання переліку торентов: %v", err)
			return
		}

		log.Printf("[GetAllTorrents] Торенти отримано")

		answer := generateAnswerList(torrents)
		helpers.SendAnswer(inputMessage.ChatId, answer)
	}

	listener := &nats_helper.NatsListener{
		Handler: processor,
	}

	nats_helper.StartNatsListener("EXECUTE_TORRENT_COMMAND_LIST", listener)
}

func generateAnswerList(torrents []transmissionrpc.Torrent) string {
	var line strings.Builder

	for _, torrent := range torrents {
		line.WriteString(
			fmt.Sprintf("%s %s\n%s %s\n%s %s\n",
				getStatusIcon(torrent), *torrent.Name,
				getProgressBar(torrent), getGigabytesLeft(torrent),
				fmt.Sprintf("%s%s%s%d", "/more_", "", "", int(*torrent.ID)),
				fmt.Sprintf("%s%s%s%d", "/files_", "", "", int(*torrent.ID)),
			),
		)
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
