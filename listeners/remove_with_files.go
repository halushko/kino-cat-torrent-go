package listeners

import (
	"context"
	"fmt"
	"github.com/hekmon/transmissionrpc/v2"
	"kino-cat-torrent-go/helpers"
	"log"
	"strconv"
)

func RemoveWithFiles() {
	processor := func(args []string, client *transmissionrpc.Client) string {
		log.Printf("[RemoveWithFiles] Старт зупинки торенту")
		strId := args[0]
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			text := fmt.Sprintf("[RemoveWithFiles] ID торента \"%s\" не валідний: %v", strId, err)
			log.Printf(text)
			return text
		}

		err = client.TorrentRemove(
			context.Background(),
			transmissionrpc.TorrentRemovePayload{
				IDs:             []int64{id},
				DeleteLocalData: true,
			})

		answer := ""
		if err != nil {
			answer = fmt.Sprintf("Торент з ID=%d не видалено", id)
		} else {
			answer = fmt.Sprintf("Торент з ID=%d видалено, файли видалено з сервера", id)
		}
		return answer
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_DELETE_WITH_FILES", processor)
}
