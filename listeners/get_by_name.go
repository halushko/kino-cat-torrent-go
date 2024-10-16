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
	processor := func(key string, args []string, client *transmissionrpc.Client) string {
		log.Printf("[GetTorrentsByName] Торенти отримано")
		torrents, err := client.TorrentGetAll(context.Background())
		if err != nil {
			text := fmt.Sprintf("[GetTorrentsByName] Помилка отримання переліку торентов: %v", err)
			log.Printf(text)
			return text
		}
		filteredTorrents := make([]transmissionrpc.Torrent, 0) // Создаем список для отфильтрованных торрентов

		searchQuery := strings.Join(args, "")
		searchQuery = strings.ReplaceAll(searchQuery, " ", "")
		searchQuery = strings.ReplaceAll(searchQuery, "_", "")
		searchQuery = strings.ReplaceAll(searchQuery, "-", "")
		searchQuery = strings.ReplaceAll(searchQuery, "/", "")
		searchQuery = strings.ReplaceAll(searchQuery, "|", "")
		searchQuery = strings.ToUpper(searchQuery)

		for i := 0; i < len(torrents); i++ {
			if torrents[i].Name != nil {
				modifiedName := strings.ReplaceAll(*torrents[i].Name, " ", "")
				modifiedName = strings.ReplaceAll(modifiedName, "_", "")
				modifiedName = strings.ReplaceAll(modifiedName, "-", "")
				modifiedName = strings.ReplaceAll(modifiedName, "/", "")
				modifiedName = strings.ReplaceAll(modifiedName, "|", "")
				modifiedName = strings.ToUpper(modifiedName)

				if strings.Contains(modifiedName, searchQuery) {
					filteredTorrents = append(filteredTorrents, torrents[i]) // Добавляем торрент в фильтрованный список
				}
			}
		}
		sort.Slice(torrents, func(i, j int) bool { return *torrents[i].ID < *torrents[j].ID })
		log.Printf("[GetTorrentsByName] Торенти для сховища %s отримано", key)
		answer := generateAnswerList(key, filteredTorrents)
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_SEARCH_BY_NAME", processor)
}
