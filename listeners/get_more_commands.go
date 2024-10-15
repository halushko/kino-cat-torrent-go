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
	processor := func(key string, args []string, client *transmissionrpc.Client) string {
		log.Printf("[GetMoreCommands] Старт отримання інформації по торенту")
		id, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			text := fmt.Sprintf("[GetMoreCommands] ID торента \"%s_%s\" не валідний: %v", args[0], args[1], err)
			log.Printf(text)
			return text
		}

		torrents, err := client.TorrentGet(context.Background(), []string{"name", "id", "status"}, []int64{id})
		if err != nil {
			text := fmt.Sprintf("[GetMoreCommands] Помилка отримання переліку торентов: %v", err)
			log.Printf(text)
			return text
		}
		answer := ""
		if len(torrents) == 1 {
			log.Printf("[GetMoreCommands] Інформацію про торент \"%d\" отримано", id)
			answer = generateAnswerMore(torrents[0], args[0], args[1])
		} else {
			log.Printf("[GetMoreCommands] Інформації про торент \"%d\" немає", id)
			answer = fmt.Sprintf("Нажаль для торента з ID=%d не можна отримати Ім'я", id)
		}
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_SHOW_COMMANDS", processor)
}

func generateAnswerMore(torrent transmissionrpc.Torrent, server string, id string) string {
	var line strings.Builder
	line.WriteString(*torrent.Name + "\n")
	line.WriteString(fmt.Sprintf("/pause_%s_%s\n", server, id))
	line.WriteString(fmt.Sprintf("/resume_%s_%s\n", server, id))
	line.WriteString(fmt.Sprintf("/info_%s_%s\n", server, id))
	line.WriteString(fmt.Sprintf("/remove_%s_%s", server, id))
	return line.String()
}
