package listeners

import (
	"context"
	"fmt"
	"github.com/hekmon/transmissionrpc/v2"
	"kino-cat-torrent-go/helpers"
	"log"
	"strconv"
)

func ExecutePauseTorrent() {
	processor := func(args []string, client *transmissionrpc.Client) string {
		log.Printf("[ExecutePauseTorrent] Старт зупинки торенту")
		strId := args[len(args)-1]
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			text := fmt.Sprintf("[ExecutePauseTorrent] ID торента \"%s\" не валідний: %v", strId, err)
			log.Printf(text)
			return text
		}

		err = client.TorrentStopIDs(context.Background(), []int64{id})

		answer := ""
		if err != nil {
			answer = fmt.Sprintf("Торент з ID=%d не зупинено", id)
		} else {
			answer = fmt.Sprintf("Торент з ID=%d зупинено", id)
		}
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_PAUSE_TORRENT", processor)
}
