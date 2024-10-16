package listeners

import (
	"context"
	"fmt"
	"github.com/hekmon/transmissionrpc/v2"
	"kino-cat-torrent-go/helpers"
	"log"
	"strconv"
	"strings"
)

func GetTorrentInfo() {
	processor := func(key string, args []string, client *transmissionrpc.Client) string {
		log.Printf("[GetTorrentInfo] Старт отримання інформації по торенту")
		strId := args[len(args)-1]
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			text := fmt.Sprintf("[GetTorrentInfo] ID торента \"%s\" не валідний: %v", strId, err)
			log.Printf(text)
			return text
		}

		answer := ""
		torrents, err := client.TorrentGet(
			context.Background(),
			[]string{"totalSize", "percentDone", "uploadedEver", "activityDate", "name", "error", "errorString", "comment", "dateCreated", "startDate"},
			[]int64{id},
		)
		if err != nil {
			answer = fmt.Sprintf("Інформацію по торенту з ID=%d не можливо отримати", id)
		} else {
			answer = generateAnswerInfo(torrents[0])
		}
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_INFO", processor)
}

func generateAnswerInfo(torrent transmissionrpc.Torrent) string {
	totalSize := torrent.TotalSize.GB()
	done := *torrent.PercentDone
	uploadedEver := *torrent.UploadedEver

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Торент %s\n/\n", *torrent.Name))
	sb.WriteString(fmt.Sprintf("Маємо: %.2f Gb", totalSize))
	sb.WriteString(fmt.Sprintf(" (%.0f%%)\n", done*100))
	//TODO sb.WriteString(fmt.Sprintf("Відвантажено: %.2f Gb", math.Round(float64(uploadedEver)/1000000.0)/1000.0))
	sb.WriteString(fmt.Sprintf(" (%.0f%%)\n", float64(uploadedEver)/totalSize))
	sb.WriteString(fmt.Sprintf("Активність: %s\n", torrent.ActivityDate.Format(`02-01-2006 15:04:05`)))
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
