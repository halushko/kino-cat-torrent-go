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

		torrents, err := client.TorrentGet(context.Background(), []string{"name", "id", "status"}, []int64{id})
		if err != nil {
			text := fmt.Sprintf("[GetMoreCommands] Помилка отримання переліку торентов: %v", err)
			log.Printf(text)
			return text
		}
		var answer string
		switch {
		case len(torrents) == 1:
			log.Printf("[GetMoreCommands] Інформацію про торент \"%d\" отримано", id)
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
	line.WriteString(fmt.Sprintf("/pause_%s\n", id))
	line.WriteString(fmt.Sprintf("/resume_%s\n", id))
	line.WriteString(fmt.Sprintf("/info_%s\n", id))
	line.WriteString(fmt.Sprintf("/remove_%s", id))
	return line.String()
}
