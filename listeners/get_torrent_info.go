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
	"strconv"
	"strings"
)

func GetTorrentInfo() {
	processor := func(msg *nats.Msg) {
		client, inputMessage := helpers.ConnectToTransmission(msg)

		log.Printf("[GetTorrentInfo] Старт зупинки торенту")
		strId := inputMessage.Arguments[len(inputMessage.Arguments)-1]
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			log.Printf("[GetTorrentInfo] ID торента \"%s\" не валідний: %v", strId, err)
			return
		}

		answer := ""
		torrents, err := client.TorrentGet(context.Background(), []string{}, []int64{id})
		if err != nil {
			answer = fmt.Sprintf("Інформацію по торенту з ID=%d не можливо отримати", id)
		} else {
			answer = generateAnswerInfo(torrents[0])
		}

		helpers.SendAnswer(inputMessage.ChatId, answer)
	}

	listener := &nats_helper.NatsListener{
		Handler: processor,
	}

	nats_helper.StartNatsListener("EXECUTE_TORRENT_COMMAND_PAUSE_TORRENT", listener)
}

func generateAnswerInfo(torrent transmissionrpc.Torrent) string {
	totalSize := float64(*torrent.TotalSize)
	done := *torrent.PercentDone
	uploadedEver := *torrent.UploadedEver
	lastActivityDate := *torrent.ActivityDate

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Торент %s\n/\n", *torrent.Name))
	sb.WriteString(fmt.Sprintf("Маємо: %f Gb", math.Round(totalSize/1000000.0)/1000.0))
	sb.WriteString(fmt.Sprintf(" (%F%%)\n", done))
	sb.WriteString(fmt.Sprintf("Відвантажено: %f Gb", math.Round(float64(uploadedEver)/1000000.0)/1000.0))
	sb.WriteString(fmt.Sprintf(" (%f%%)\n", (10*float64(uploadedEver)/totalSize)/10.0))
	sb.WriteString(fmt.Sprintf("Активність: %s\n", lastActivityDate))
	if *torrent.Error != 0 {
		sb.WriteString(fmt.Sprintf("Помилка: %s\n", *torrent.ErrorString))
	}
	sb.WriteString(fmt.Sprintf("Торент створено: %s\n", torrent.DateCreated.Format(`02-01-2006 15:04:05`)))
	sb.WriteString(fmt.Sprintf("Початок закачки: %s\n", torrent.StartDate.Format("02-01-2006 15:04:05")))
	if *torrent.Comment != "" {
		sb.WriteString(fmt.Sprintf("Інфа: %s\n", *torrent.Comment))
	}
	return sb.String()
}
