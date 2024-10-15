package listeners

import (
	"context"
	"fmt"
	"github.com/hekmon/transmissionrpc/v2"
	"kino-cat-torrent-go/helpers"
	"log"
	"strconv"
)

func GetMoreCommands() {
	processor := func(key string, args []string, client *transmissionrpc.Client) string {
		log.Printf("[GetMoreCommands] Старт отримання інформації по торенту")
		strId := args[len(args)-1]
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			text := fmt.Sprintf("[GetMoreCommands] ID торента \"%s\" не валідний: %v", strId, err)
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
			answer = generateAnswerMore(torrents[0])
		} else {
			log.Printf("[GetMoreCommands] Інформації про торент \"%d\" немає", id)
			answer = fmt.Sprintf("Нажаль для торента з ID=%d не можна отримати Ім'я", id)
		}
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_SHOW_COMMANDS", processor)
}

func generateAnswerMore(torrent transmissionrpc.Torrent) string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
		*torrent.Name,
		fmt.Sprintf("%s%s%s%d", "/pause_", "", "", *torrent.ID),
		fmt.Sprintf("%s%s%s%d", "/resume_", "", "", *torrent.ID),
		fmt.Sprintf("%s%s%s%d", "/info_", "", "", *torrent.ID),
		fmt.Sprintf("%s%s%s%d", "/remove_", "", "", *torrent.ID),
	)
}
