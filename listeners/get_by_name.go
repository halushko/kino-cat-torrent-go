package listeners

import (
	"context"
	"fmt"
	"github.com/hekmon/transmissionrpc/v2"
	"kino-cat-torrent-go/helpers"
	"log"
	"sort"
	"strings"
)

func GetTorrentsByName() {
	processor := func(args []string, client *transmissionrpc.Client) []string {
		log.Printf("[GetTorrentsByName] Торенти отримано")
		torrents, err := client.TorrentGetAll(context.Background())
		if err != nil {
			text := fmt.Sprintf("[GetTorrentsByName] Помилка отримання переліку торентов: %v", err)
			log.Printf(text)
			return []string{text}
		}
		filteredTorrents := make([]transmissionrpc.Torrent, 0)

		searchQuery := prepareText(strings.Join(args, ""))

		for i := 0; i < len(torrents); i++ {
			if torrents[i].Name != nil {
				modifiedName := prepareText(strings.ReplaceAll(*torrents[i].Name, " ", ""))

				if strings.Contains(modifiedName, searchQuery) {
					filteredTorrents = append(filteredTorrents, torrents[i])
				}
			}
		}
		sort.Slice(torrents, func(i, j int) bool { return *torrents[i].ID < *torrents[j].ID })
		log.Printf("[GetTorrentsByName] Торенти отримано")
		var answer []string
		switch {
		case len(filteredTorrents) > 0:
			answer = generateAnswerList(filteredTorrents)
		default:
			answer = []string{"Нажаль торента з таким ім'ям не знайдено"}
		}
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_SEARCH_BY_NAME", processor)
}

func prepareText(text string) string {
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "_", "")
	text = strings.ReplaceAll(text, "-", "")
	text = strings.ReplaceAll(text, "/", "")
	text = strings.ReplaceAll(text, "|", "")
	text = strings.ToUpper(text)
	return text
}
