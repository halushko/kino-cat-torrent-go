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

func AskDeleteTorrent() {
	processor := func(args []string, client *transmissionrpc.Client) string {
		log.Printf("[AskDeleteTorrent] Старт генерування форми підтвердження")
		strId := args[0]
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			text := fmt.Sprintf("[AskDeleteTorrent] ID торента \"%s\" не валідний: %v", strId, err)
			log.Printf(text)
			return text
		}

		torrents, err := client.TorrentGet(context.Background(), []string{"name", "id", "status"}, []int64{id})
		if err != nil {
			text := fmt.Sprintf("[AskDeleteTorrent] Помилка отримання переліку торентов: %v", err)
			log.Printf(text)
			return text
		}
		var answer string
		switch {
		case len(torrents) == 1:
			log.Printf("[AskDeleteTorrent] Інформацію про торент \"%d\" отримано", id)
			answer = generateAnswerAskDelete(torrents[0], id)
		default:
			log.Printf("[AskDeleteTorrent] Інформації про торент \"%d\" немає", id)
			answer = fmt.Sprintf("Нажаль для торента з ID=%d не можна отримати Ім'я", id)
		}
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_DELETE", processor)
}

func generateAnswerAskDelete(torrent transmissionrpc.Torrent, id int64) string {
	var line strings.Builder
	line.WriteString(fmt.Sprintf("Дійсно хочете видалити торент %s\n", *torrent.Name))
	line.WriteString(fmt.Sprintf("-- видалити лише сам торент: /approve_just_torrent_%d\n", id))
	line.WriteString(fmt.Sprintf("-- видалити ще й скачані файли: /approve_with_files_%d", id))
	return line.String()
}
