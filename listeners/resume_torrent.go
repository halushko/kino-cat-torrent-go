package listeners

import (
	"context"
	"fmt"
	"github.com/hekmon/transmissionrpc/v2"
	"kino-cat-torrent-go/helpers"
	"log"
	"strconv"
)

func ExecuteResumeTorrent() {
	processor := func(args []string, client *transmissionrpc.Client) string {
		log.Printf("[ExecuteResumeTorrent] Старт поновлення торенту")
		strId := args[0]
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			text := fmt.Sprintf("[ExecuteResumeTorrent] ID торента \"%s\" не валідний: %v", strId, err)
			log.Printf(text)
			return text
		}

		err = client.TorrentStartIDs(context.Background(), []int64{id})

		var answer string
		switch {
		case err != nil:
			answer = fmt.Sprintf("Торент з ID=%d не поновлено", id)
		default:
			answer = fmt.Sprintf("Торент з ID=%d поновлено", id)
		}
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_RESUME_TORRENT", processor)
}
