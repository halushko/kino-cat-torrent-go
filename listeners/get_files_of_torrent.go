package listeners

import (
	"context"
	"fmt"
	"github.com/hekmon/transmissionrpc/v2"
	"kino-cat-torrent-go/helpers"
	"log"
	"sort"
	"strconv"
	"strings"
)

func GetFilesOfTorrent() {
	processor := func(args []string, client *transmissionrpc.Client) string {
		log.Printf("[GetFilesOfTorrent] Старт відображення файлів")
		strId := args[0]
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			text := fmt.Sprintf("[GetFilesOfTorrent] ID торента \"%s\" не валідний: %v", strId, err)
			log.Printf(text)
			return text
		}

		torrents, err := client.TorrentGet(
			context.Background(),
			[]string{"name", "files"},
			[]int64{id},
		)

		var answer string
		switch {
		case err != nil:
			answer = fmt.Sprintf("Файли торента з ID=%d не знайдено", id)
		default:
			answer = getInfoAboutFiles(torrents[0])
		}
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_LIST_FILES", processor)
}

func getInfoAboutFiles(torrent transmissionrpc.Torrent) string {
	files := torrent.Files
	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s\n/\n", *torrent.Name))

	for _, file := range files {
		done := float64(file.BytesCompleted) / float64(file.Length)
		pb := getProgressBar(done, 10)
		name := file.Name

		if strings.HasPrefix(name, *torrent.Name+"/") {
			name = strings.TrimPrefix(name, *torrent.Name+"/")
		}

		percent := ""
		if done >= 1 {
			percent = "заверш"
		} else {
			percent = fmt.Sprintf("%.2f Gb", float64(file.Length-file.BytesCompleted)/1024.0/1024.0/1024.0)
		}
		sb.WriteString(fmt.Sprintf("%s\n", name))
		sb.WriteString(fmt.Sprintf("%s (%s)\n", pb, percent))
	}
	return sb.String()
}
