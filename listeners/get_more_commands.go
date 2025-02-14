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

func GetMoreCommands() {
	processor := func(args []string, client *transmissionrpc.Client) string {
		log.Printf("[GetMoreCommands] Старт отримання інформації по торенту")
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			text := fmt.Sprintf("[GetMoreCommands] ID торента \"%s\" не валідний: %v", args[0], err)
			log.Printf(text)
			return text
		}

		torrents, err := client.TorrentGet(context.Background(), []string{"name", "id", "status", "downloadDir"}, []int64{id})
		if err != nil {
			text := fmt.Sprintf("[GetMoreCommands] Помилка отримання переліку торентов: %v", err)
			log.Printf(text)
			return text
		}
		var answer string
		switch {
		case len(torrents) == 1:
			log.Printf("[GetMoreCommands] Інформацію про торент \"%d\" отримано", id)
			//log.Printf("[GetMoreCommands] торент \"%v\" отримано", torrents[0])
			answer = generateAnswerMore(torrents[0], args[0])
		default:
			log.Printf("[GetMoreCommands] Інформації про торент \"%d\" немає", id)
			answer = fmt.Sprintf("Нажаль для торента з ID=%d не можна отримати Ім'я", id)
		}
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_SHOW_COMMANDS", processor)
}

func generateAnswerMore(torrent transmissionrpc.Torrent, id string) string {
	var line strings.Builder
	line.WriteString(*torrent.Name + "\n")

	line.WriteString(fmt.Sprintf("/info_%s\n", id))
	if torrent.Status != nil && *torrent.Status == transmissionrpc.TorrentStatusStopped {
		line.WriteString(fmt.Sprintf("/resume_%s\n", id))
	} else {
		line.WriteString(fmt.Sprintf("/pause_%s\n", id))
	}
	line.WriteString(fmt.Sprintf("/remove_%s\n", id))

	if torrent.DownloadDir != nil {
		log.Printf("[generateAnswerMore] torrent.DownloadDir == %s", *torrent.DownloadDir)
	}
	if torrent.DownloadDir != nil && *torrent.DownloadDir == helpers.DownloadDir {
		line.WriteString(fmt.Sprintf("/backlog_%s\n", id))
	} else {
		line.WriteString(fmt.Sprintf("/de_backlog_%s\n", id))
	}

	return line.String()
}
