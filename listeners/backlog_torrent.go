package listeners

import (
	"context"
	"fmt"
	"github.com/hekmon/transmissionrpc/v2"
	"kino-cat-torrent-go/helpers"
	"log"
	"strconv"
)

func ExecuteBacklogTorrent() {
	processor := func(args []string, client *transmissionrpc.Client) string {
		log.Printf("[ExecuteBacklogTorrent] Старт перенесення торенту у сховище для подивитися пізніше")

		return MoveTorrent(args, "/downloads/backlog", client)
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_BACKLOG", processor)
}

func ExecuteDeBacklogTorrent() {
	processor := func(args []string, client *transmissionrpc.Client) string {
		log.Printf("[ExecuteDeBacklogTorrent] Старт перенесення торенту зі сховище для подивитися пізніше")

		return MoveTorrent(args, "/downloads/complete", client)
	}

	helpers.ListenToNatsMessages("EXECUTE_TORRENT_COMMAND_DEBACKLOG", processor)
}

func MoveTorrent(args []string, newLocation string, client *transmissionrpc.Client) string {
	strId := args[0]
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		text := fmt.Sprintf("[MoveTorrent] ID торента \"%s\" не валідний: %v", strId, err)
		log.Printf(text)
		return text
	}

	err = client.TorrentSetLocation(context.Background(), id, newLocation, true)

	var answer string
	switch {
	case err != nil:
		answer = fmt.Sprintf("[MoveTorrent] Помилка при зміні локації медіафайлів торента ID=\"%d\": %v", id, err)
	default:
		answer = fmt.Sprintf("[MoveTorrent] Почато операцію переміщення торанта ID=\"%d\"", id)
	}
	log.Printf(answer)
	return answer
}
