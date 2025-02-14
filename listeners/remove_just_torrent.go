package listeners

import (
	"context"
	"fmt"
	"github.com/hekmon/transmissionrpc/v2"
	"kino-cat-torrent-go/helpers"
	"log"
	"strconv"
)

func RemoveJustTorrent() {
	processor := func(args []string, client *transmissionrpc.Client) string {
		log.Printf("[RemoveJustTorrent] Старт зупинки торенту")
		strId := args[0]
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			text := fmt.Sprintf("[RemoveJustTorrent] ID торента \"%s\" не валідний: %v", strId, err)
			log.Printf(text)
			return text
		}

		err = client.TorrentRemove(
			context.Background(),
			transmissionrpc.TorrentRemovePayload{
				IDs:             []int64{id},
				DeleteLocalData: false,
			})

		var answer string
		switch {
		case err != nil:
			answer = fmt.Sprintf("Торент з ID=%d не видалено", id)
		default:
			answer = fmt.Sprintf("Торент з ID=%d видалено, файли залишилися на сервері", id)
		}
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_DELETE_ONLY_TORRENT", processor)
}
